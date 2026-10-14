package main

import (
	"log"
	"my_work/internal/config"
	"my_work/internal/database"
	"my_work/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var cfg *config.Config
	var err error
	cfg, err = config.Load()

	if err != nil {
		log.Fatal("Failed to load configuration: ", err)
	}
	var pool *pgxpool.Pool
	pool, err = database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
		return
	}
	defer pool.Close()

	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "my work API is running",
			"status":   "success",
			"database": "connected",
		})
	})

	router.POST("/auth/register", handlers.CreateMemberHandler(pool))

	router.POST("/myWorks", handlers.CreateMyWorkHandler(pool))
	router.GET("/myWorks", handlers.GetAllMyWorksHandler(pool))
	router.GET("/myWorks/:id", handlers.GetMyWorkByIDHandler(pool))
	router.PUT("/myWorks/:id", handlers.UpdateMyWorkHandler(pool))
	router.DELETE("/myWorks/:id", handlers.DeleteMyWorkHandler(pool))

	router.Run(":" + cfg.Port)
}
