package serv

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"tgtest/domain"
	"time"
	"unicode"

	"github.com/segmentio/kafka-go"
	"golang.org/x/crypto/bcrypt"
)

type EventService struct {
	writer *kafka.Writer
}
type RepoInterface interface {
	CreateUser(ctx context.Context, user *domain.User) (string, error)
}

type Service struct {
	repo RepoInterface
}

func CreateEventService(writer *kafka.Writer) *EventService {
	return &EventService{writer: writer}
}

type UserCreatedEvent struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateServ(repo RepoInterface) *Service {
	return &Service{repo: repo}
}

func (e *EventService) PublishEvent(ctx context.Context, topic string, event *UserCreatedEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	e.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.ID),
		Topic: topic,
		Value: data,
	})
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
	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}
	return id, nil
}
