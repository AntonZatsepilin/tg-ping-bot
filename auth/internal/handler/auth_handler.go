package handler

import (
	"context"
	"goPingRobot/auth/internal/service"
	"goPingRobot/auth/proto"
)

type AuthHandler struct {
    proto.UnimplementedAuthServiceServer
    service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
    return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
    err := h.service.Register(req.Username, req.Password)
    if err != nil {
        return nil, err
    }
    return &proto.RegisterResponse{Message: "User registered successfully"}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
    token, err := h.service.Login(req.Username, req.Password)
    if err != nil {
        return nil, err
    }
    return &proto.LoginResponse{Token: token}, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *proto.ValidateTokenRequest) (*proto.ValidateTokenResponse, error) {
    userID, valid := h.service.ValidateToken(req.Token)
    return &proto.ValidateTokenResponse{UserId: int32(userID), Valid: valid}, nil
}
