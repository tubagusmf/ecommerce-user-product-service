package model

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ContextAuthKey string

const BearerAuthKey ContextAuthKey = "BearerAuth"

type IUserRepository interface {
	FindAll(ctx context.Context, user User) ([]*User, error)
	FindById(ctx context.Context, id int64) (*User, error)
	FindByEmail(ctx context.Context, email string) *User
	Create(ctx context.Context, user User) (*User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id int64) error
	CreateSession(ctx context.Context, session UserSession) (*UserSession, error)
	FindSessionByToken(ctx context.Context, token string) (*UserSession, error)
	DeleteSession(ctx context.Context, token string) error
}

type IUserUsecase interface {
	FindAll(ctx context.Context, user User) ([]*User, error)
	FindById(ctx context.Context, id int64) (*User, error)
	Create(ctx context.Context, in CreateUserInput) (token string, err error)
	Update(ctx context.Context, id int64, in UpdateUserInput) error
	Delete(ctx context.Context, id int64) error
	ValidateSession(ctx context.Context, token string) (*UserSession, error)
	Login(ctx context.Context, in LoginInput) (token string, err error)
	Logout(ctx context.Context, token string) error
}

type CustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

type User struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

type UserSession struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserInput struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3"`
	Role     string `json:"role" validate:"required"`
}

type UpdateUserInput struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3"`
	Role     string `json:"role" validate:"required"`
}
