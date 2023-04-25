package services

import (
	"auth_user/app/db"
	"auth_user/app/models"
	"auth_user/app/pb"
	"auth_user/app/utils"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	H   db.Handler
	Jwt utils.JwtWrapper
}

func (s *Server) StartHttp(rg *gin.RouterGroup) {
	rg.POST("/register", s.registerUser)
	rg.POST("/login", s.loginUser)
	rg.POST("/validate", s.validateUser)

}

func (s *Server) validateUser(ctx *gin.Context) {
	const BEARER_SCHEMA = "Bearer "
	authHeader := ctx.GetHeader("Authorization")
	fmt.Println("auth:", authHeader)
	token := authHeader[len(BEARER_SCHEMA):]
	val, errs := s.Validate(ctx, &pb.ValidateRequest{Token: token})
	if errs != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": errs.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"user_id": val.UserId})

	return

}

func (s *Server) registerUser(ctx *gin.Context) {
	var req RegisterReq
	var user models.User
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	if result := s.H.DB.Where(&models.User{Email: req.Email}).First(&user); result.Error == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": "email_already-exist"})
		return
	}

	user.Email = req.Email
	user.Password = utils.HashPassword(req.Password)

	s.H.DB.Create(&user)
}

func (s *Server) loginUser(ctx *gin.Context) {
	var req RegisterReq
	var user models.User
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	if result := s.H.DB.Where(&models.User{Email: req.Email}).First(&user); result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": "not-found"})
		return
	}

	match := utils.CheckPasswordHash(req.Password, user.Password)

	if !match {
		ctx.JSON(http.StatusNotFound, gin.H{"code": "not-found"})
		return
	}

	token, _ := s.Jwt.GenerateToken(user)

	ctx.JSON(http.StatusOK, gin.H{"token": token, "userId": user.Id})
	return

}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var user models.User

	if result := s.H.DB.Where(&models.User{Email: req.Email}).First(&user); result.Error == nil {
		return &pb.RegisterResponse{
			Status: http.StatusConflict,
			Error:  "E-Mail already exists",
		}, nil
	}

	user.Email = req.Email
	user.Password = utils.HashPassword(req.Password)

	s.H.DB.Create(&user)

	return &pb.RegisterResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var user models.User

	if result := s.H.DB.Where(&models.User{Email: req.Email}).First(&user); result.Error != nil {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	match := utils.CheckPasswordHash(req.Password, user.Password)

	if !match {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	token, _ := s.Jwt.GenerateToken(user)

	return &pb.LoginResponse{
		Status: http.StatusOK,
		Token:  token,
	}, nil
}

func (s *Server) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	fmt.Println("in valll", req.Token)
	claims, err := s.Jwt.ValidateToken(req.Token)
	fmt.Println(err, claims, "><")
	if err != nil {
		return &pb.ValidateResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	var user models.User

	if result := s.H.DB.Where(&models.User{Email: claims.Email}).First(&user); result.Error != nil {
		return &pb.ValidateResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	return &pb.ValidateResponse{
		Status: http.StatusOK,
		UserId: user.Id,
	}, nil
}
