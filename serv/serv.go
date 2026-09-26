package serv

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"tgtest/apiclient"
	"tgtest/domain"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EventService struct {
	writer *kafka.Writer
}
type RepoInterface interface {
	CreateUser(ctx context.Context, topic string, user *domain.User) (string, error)
	UpdateBalance(ctx context.Context, userID string, amount int64) error
}

func CreateEventService(writer *kafka.Writer) *EventService {
	return &EventService{writer: writer}
}

type UserCreatedEvent struct {
	UserID    string    `json:"user_id"`
	EventID   string    `json:"event_id"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"created_at"`
}
type Service struct {
	repo   RepoInterface
	client apiclient.BalanceClientRPC
}

func CreateServ(repo RepoInterface) (*Service, error) {
	client, err := apiclient.CreateClientRPC()
	if err != nil {
		return nil, err
	}
	return &Service{client: client, repo: repo}, nil
}

func (e *Service) GetBalance(ctx context.Context, userID string) (string, error) {
	response, err := e.client.GetBalance(ctx, userID)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return "", err
		}
		switch st.Code() {
		case codes.NotFound:
			return "", domain.ErrUserNotFound
		case codes.DeadlineExceeded:
			return "", context.DeadlineExceeded
		case codes.InvalidArgument:
			return "", domain.ErrInvalidInput
		default:
			return "", fmt.Errorf("get balance failed: %w", err)

		}
	}
	if err := e.repo.UpdateBalance(ctx, userID, response.Balance); err != nil {
		return "", err
	}
	balance := strconv.Itoa(int(response.Balance))
	return balance, nil
}

func (e *EventService) PublishEvent(ctx context.Context, topic string, event *UserCreatedEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = e.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.UserID),
		Topic: topic,
		Value: data,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) CreateUser(ctx context.Context, name, email, password, age string) (string, error) {
	ageInt, err := strconv.Atoi(age)
	if err != nil {
		return "", fmt.Errorf("failed to convert age")
	}
	if ageInt < 18 {
		return "", domain.ErrAgeForbidden
	}
	if len(password) < 8 {
		return "", domain.ErrShortPassword
	}
	var hasUpper, hasSpecial bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsPunct(ch), unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}
	if !hasUpper {
		return "", domain.ErrNoUpperLetter
	}
	if !hasSpecial {
		return "", domain.ErrNoSpecialChar
	}
	name = strings.TrimSpace(name)

	if name == "" || len(name) < 2 || len(name) > 70 {
		return "", domain.ErrInvalidName
	}
	if !strings.Contains(email, "@") || len(email) > 70 {
		return "", domain.ErrInvalidEmail
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt password, %w", err)
	}
	var user = &domain.User{
		Name:     name,
		Email:    email,
		Password: string(pass),
		Age:      age,
	}
	eventID := uuid.New()
	id, err := s.repo.CreateUser(ctx, eventID.String(), user)
	if err != nil {
		return "", err
	}
	return id, nil
}
