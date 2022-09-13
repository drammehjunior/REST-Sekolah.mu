package api

import (
	"exampleclean.com/refactor/app/api/handler"
	"github.com/gin-gonic/gin"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(userHandler *handler.UserHandler) *ServerHTTP {
	engine := gin.New()

	// Use logger from Gin
	engine.Use(gin.Logger())

	//// Swagger docs
	//engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	//
	//// Request JWT
	engine.POST("login", userHandler.LoginHandler)
	engine.PUT("user/updatePassword", userHandler.UpdatePassword)

	// Auth middleware
	api := engine.Group("/api")
	api.GET("users", userHandler.FindAll)
	api.GET("users/:id", userHandler.FindByID)
	api.POST("users", userHandler.SaveSignup)
	api.DELETE("users/:id", userHandler.Delete)
	api.GET("user/email/:mail", userHandler.FindByEmail)
	//api.POST("login", userHandler.UserLogin)

	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() {
	sh.engine.Run(":8080")
}
