package handlers

import (
	"my_work/internal/config"
	"my_work/internal/models"
	"my_work/internal/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func CreateMemberHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var registerRequest RegisterRequest

		if err := c.BindJSON(&registerRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(registerRequest.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters long!"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password" + err.Error()})
			return
		}

		member := &models.Member{
			Email:    registerRequest.Email,
			Password: string(hashedPassword),
		}
		createdMember, err := repository.CreateMember(pool, member)
		if err != nil {
			if err.Error() != "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, createdMember)
	}
}

func LoginHandler(pool *pgxpool.Pool, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest LoginRequest
		if err := c.BindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		member, err := repository.GetMemberByEmail(pool, loginRequest.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(loginRequest.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Credentials"})
			return
		}
		//map[string]any{}
		claims := jwt.MapClaims{
			"member_id": member.ID,
			"email":     member.Email,
			"exp":       time.Now().Add(24 * time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(cfg.JWTSecret)) //generate signature
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, LoginResponse{Token: tokenString})

	}

}

func TestProtectedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		memberID, exists := c.Get("member_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "member_id not found in context"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":   "Protected route accessed successfully",
			"member_id": memberID,
		})

	}
}
