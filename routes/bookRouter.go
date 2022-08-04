package routes

import (
	controllers "github.com/kapbyte/book-app-tutorial/controller"

	"github.com/gin-gonic/gin"
)

func BookRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("book/create", controllers.CreateBook())
	incomingRoutes.GET("book/:book_id", controllers.GetBook())
	incomingRoutes.PATCH("book/:book_id", controllers.UpdateBook())
}
