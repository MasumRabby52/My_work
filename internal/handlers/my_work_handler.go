package handlers

import (
	"my_work/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateMyWorkInput struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}
type UpdateMyWorkInput struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

func CreateMyWorkHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		memberIDInterface, exists := c.Get("member_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "member_id not found in context"})
			return
		}
		memberID := memberIDInterface.(string)

		var input CreateMyWorkInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		myWork, err := repository.CreateMyWork(pool, input.Title, input.Completed, memberID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, myWork)
	}
}

func GetAllMyWorksHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		memberIDInterface, exists := c.Get("member_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "member_id not found in context"})
			return
		}
		memberID := memberIDInterface.(string)

		myWorks, err := repository.GetAllMyWorks(pool, memberID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, myWorks)
	}
}

func GetMyWorkByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		memberIDInterface, exists := c.Get("member_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "member_id not found in context"})
			return
		}
		memberID := memberIDInterface.(string)

		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid myWork ID"})
			return
		}
		myWork, err := repository.GetMyWorkByID(pool, id, memberID)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "myWork Not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
			return
		}
		c.JSON(http.StatusOK, myWork)
	}
}

func UpdateMyWorkHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		memberIDInterface, exists := c.Get("member_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "member_id not found in context"})
			return
		}
		memberID := memberIDInterface.(string)

		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid myWork ID"})
		}

		var input UpdateMyWorkInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
			return
		}

		if input.Title == nil && input.Completed == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "At least one field(title or completed) must be provided"})
			return
		}

		existing, err := repository.GetMyWorkByID(pool, id, memberID)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "myWork not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		title := existing.Title
		if input.Title != nil {
			title = *input.Title
		}

		completed := existing.Completed
		if input.Completed != nil {
			completed = *input.Completed
		}
		myWork, err := repository.UpdateMyWork(pool, id, title, completed, memberID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, myWork)
	}
}

func DeleteMyWorkHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		memberIDInterface, exists := c.Get("member_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "member_id not found in context"})
			return
		}
		memberID := memberIDInterface.(string)

		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid myWork ID"})
		}
		err = repository.DeleteMyWork(pool, id, memberID)
		if err != nil {
			if err.Error() == "myWork with id "+idStr+" not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": "myWork not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "myWork deleted successfully"})
	}
}
