package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine{
	r:=gin.Default()
	r.Static("/uploads","./uploads")
	r.POST("/register",controllers.Register)
	r.POST("/login",controllers.Login)
	api:=r.Group("/api")
	{
		api.GET("post",controllers.GetMovie)
		
		protected:=api.Group("/",middlewares.AuthMiddleware())
		{
		protected.GET("post/:id",controllers.GetByID)
		protected.POST("post",controllers.CreateMovie)
		protected.PUT("/post/:id",controllers.UpdateMovie)
		protected.DELETE("/post/:id",controllers.DeleteMovie)
		}
	}

return  r	
}