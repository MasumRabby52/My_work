package main

import "github.com/gin-gonic/gin"

func main() {
	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "my work API is running",
			"status":   "success",
			"database": "connected",
		})
	})
	router.Run(":8000")
}
