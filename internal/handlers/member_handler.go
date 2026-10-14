package handlers

import (
	"my_work/internal/models"
	"my_work/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
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
