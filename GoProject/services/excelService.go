package services

import (
	"GoProject/config"
	"GoProject/models"
	"GoProject/utils"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func HandleFileUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "File upload failed")
		return
	}

	savePath := filepath.Join("uploads", filepath.Base(file.Filename))
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to save file")
		return
	}

	go func() {
		if err := parseExcelAsync(savePath); err != nil {
			log.Println("Error during Excel parsing:", err)
		}
	}()

	utils.RespondSuccess(c, http.StatusAccepted, "File uploaded. Processing in background.", nil)
}

func parseExcelAsync(path string) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return errors.New("failed to open Excel file")
	}
	defer f.Close()

	sheet := f.GetSheetName(0)
	if sheet == "" {
		return errors.New("no sheets found in Excel file")
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		return errors.New("failed to read rows from sheet")
	}

	if len(rows) < 1 {
		return errors.New("empty Excel file")
	}

	var records []models.Record
	for _, row := range rows[1:] {
		if len(row) < 10 {
			continue
		}
		record := models.Record{
			FirstName:   row[0],
			LastName:    row[1],
			CompanyName: row[2],
			Address:     row[3],
			City:        row[4],
			County:      row[5],
			Postal:      row[6],
			Phone:       row[7],
			Email:       row[8],
			Web:         row[9],
		}
		records = append(records, record)
	}

	if len(records) == 0 {
		return errors.New("no valid records to cache")
	}

	data, err := json.Marshal(records)
	if err != nil {
		return errors.New("failed to marshal records for cache")
	}
	err = config.ConnectDB.Create(&records).Error

	if err != nil {
		return errors.New("failed to marshal records for cache")
	}

	if err := config.RedisClient.Set(context.Background(), "imported_data", data, 5*time.Minute).Err(); err != nil {
		return errors.New("failed to set cache in Redis")
	}

	log.Println("Import complete:", len(records), "records processed and cached.")
	return nil
}

func FetchRecords(c *gin.Context) {
	ctx := context.Background()
	cached, err := config.RedisClient.Get(ctx, "imported_data").Result()
	if err == nil {
		var recs []models.Record
		if json.Unmarshal([]byte(cached), &recs) == nil {
			utils.RespondSuccess(c, http.StatusOK, "Records fetched from cache", recs)
			return
		}
	}
	utils.RespondError(c, http.StatusNotFound, "No cached data available")
}

func UpdateRecordByEmail(c *gin.Context) {
	email := c.Param("email")
	var updated models.Record
	if err := c.BindJSON(&updated); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid JSON")
		return
	}

	ctx := context.Background()
	cached, err := config.RedisClient.Get(ctx, "imported_data").Result()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to get cache")
		return
	}

	var records []models.Record
	if err := json.Unmarshal([]byte(cached), &records); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to parse cached records")
		return
	}

	updatedFlag := false
	for i, rec := range records {
		if rec.Email == email {
			records[i] = updated
			updatedFlag = true
			break
		}
	}

	if !updatedFlag {
		utils.RespondError(c, http.StatusNotFound, "Record not found")
		return
	}

	data, err := json.Marshal(records)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to marshal updated records")
		return
	}

	if err := config.RedisClient.Set(ctx, "imported_data", data, 5*time.Minute).Err(); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to update cache")
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "Record updated successfully", nil)
}

