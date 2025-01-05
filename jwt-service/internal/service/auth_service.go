package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"goPingRobot/auth/internal/models"
	"goPingRobot/auth/internal/repository"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/argon2"
)

type AuthService struct {
    repo      *repository.UserRepository
    jwtSecret string
}

func NewAuthService(repo *repository.UserRepository, jwtSecret string) *AuthService {
    return &AuthService{repo: repo, jwtSecret: jwtSecret}
}

func (s *AuthService) Register(username, password string) error {
	salt, err := generateSalt(16)
    if err != nil {
        return err
    }
    hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	
    user := &models.User{Username: username, PasswordHash: string(hash)}
    return s.repo.Create(user)
}

func (s *AuthService) Login(username, password string) (string, error) {
    user, err := s.repo.GetByUsername(username)
    if err != nil {
        return "", err
    }

    parts := strings.Split(user.PasswordHash, ":")
    if len(parts) != 2 {
        return "", errors.New("invalid stored password hash")
    }

    salt, err := base64.RawStdEncoding.DecodeString(parts[0])
    if err != nil {
        return "", errors.New("invalid stored password hash")
    }

    storedHash, err := base64.RawStdEncoding.DecodeString(parts[1])
    if err != nil {
        return "", errors.New("invalid stored password hash")
    }

    hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

    if !compareHashes(storedHash, hash) {
        return "", errors.New("invalid credentials")
	}

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    })
    return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateToken(tokenStr string) (int, bool) {
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        return []byte(s.jwtSecret), nil
    })
    if err != nil || !token.Valid {
        return 0, false
    }
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return 0, false
    }
    userID, ok := claims["user_id"].(float64)
    return int(userID), ok
}

func generateSalt(size int) ([]byte, error) {
    salt := make([]byte, size)
    _, err := rand.Read(salt)
    if err != nil {
        return nil, err
    }
    return salt, nil
}


func compareHashes(a, b []byte) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}