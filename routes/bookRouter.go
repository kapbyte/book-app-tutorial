package routes

import (
	controllers "github.com/kapbyte/book-app-tutorial/controller"

	"github.com/gin-gonic/gin"
)

func BookRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("books/create", controllers.CreateBook())
	incomingRoutes.GET("books/:book_id", controllers.GetBook())
	incomingRoutes.PATCH("books/:book_id", controllers.UpdateBook())
	incomingRoutes.DELETE("books/:book_id", controllers.DeleteBook())
	incomingRoutes.GET("books", controllers.GetAllBooks())
}
