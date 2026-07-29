package user

import (
	"context"
	"errors"
	"time"

	"github.com/TNJKL/bookmark-management/pkg/dbutils"
	"github.com/golang-jwt/jwt/v5"
)

const tokenDuration = 24 * time.Hour

// ErrInvalidCredentials is returned when the username or password does not match
var ErrInvalidCredentials = errors.New("invalid credentials")

// Login authenticates a user and returns a signed JWT token if successful
func (s *service) Login(ctx context.Context, username, password string) (string, error) {
	//check user exists with username
	user, err := s.repo.GetUserByUsername(ctx, username)
	switch {
	case errors.Is(err, dbutils.ErrRecordNotFound):
		return "", ErrInvalidCredentials
	case err == nil:
	default:
		return "", err
	}
	//compare password hash
	if !s.hasher.Compare(user.Password, password) {
		return "", ErrInvalidCredentials
	}
	//if match -- generate token
	tokenContent := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(tokenDuration).Unix(),
	}

	tokenString, err := s.jwtGenerator.GenerateJWT(tokenContent)
	if err != nil {
		return "", err
	}
	//return token
	return tokenString, nil
}
