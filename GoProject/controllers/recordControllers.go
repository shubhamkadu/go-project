package controllers

import (
	"GoProject/services"
	"github.com/gin-gonic/gin"
)

func UploadExcel(c *gin.Context) {
	services.HandleFileUpload(c)
}

func ViewRecords(c *gin.Context) {
	services.FetchRecords(c)
}

func UpdateRecord(c *gin.Context) {
	services.UpdateRecordByEmail(c)
}
