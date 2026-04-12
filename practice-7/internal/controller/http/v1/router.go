package v1

import (
	"practice-7/internal/usecase"
	"practice-7/pkg/logger"
	"practice-7/utils"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, uc usecase.UserInterface, l logger.Interface) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	handler.Use(utils.RateLimitMiddleware()) // Problem 3: applied globally

	v1 := handler.Group("/v1")
	{
		newUserRoutes(v1, uc, l)
	}
}
