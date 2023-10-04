package user

import (
	"context"
	"crypto/sha1"
	"fmt"

	"github.com/MateSousa/api/internal/telemetry"
	"go.opentelemetry.io/otel/codes"
)

type UseCase interface {
	ValidateUser(ctx context.Context, email, password string) error
}

type Service struct {
	repo      Repository
	telemetry telemetry.Telemetry
}

func NewService(repo Repository, telemetry telemetry.Telemetry) *Service {
	return &Service{repo: repo, telemetry: telemetry}
}

func (s *Service) ValidateUser(ctx context.Context, email, password string) error {
	span := s.telemetry.Tracer().Start(ctx, "ValidateUser")
	defer span.End()

	u, err := s.repo.Get(ctx, email)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	if u == nil {
		err := fmt.Errorf("user not found")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return s.ValidatePassword(ctx, u, password)
}

func (s *Service) ValidatePassword(ctx context.Context, u *User, password string) error {
	span := s.telemetry.Tracer().Start(ctx, "ValidatePassword")
	defer span.End()

	h := sha1.New()
	h.Write([]byte(password))

	p := fmt.Sprintf("%x", h.Sum(nil))
	if p != u.Password {
		err := fmt.Errorf("invalid password")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}
