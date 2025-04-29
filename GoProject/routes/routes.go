package routes

import (
	"GoProject/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/upload", controllers.UploadExcel)
	r.GET("/records", controllers.ViewRecords)
	r.PUT("/records/:email", controllers.UpdateRecord)
}
