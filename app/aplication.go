package app

import (
	"UdemyApp/bookstore_users-api/logger"
	"github.com/gin-gonic/gin"
)

var (    
	router = gin.Default()
)

func StartApplication() {
	mapUrls()
	logger.Info("about to start aplication...")
	router.Run(":8080")
}