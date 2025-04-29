package services

import (
	"GoProject/config"
	"GoProject/models"
	"context"
	"encoding/json"
	"errors"
	"time"
)

// CacheRecords caches the records in Redis
func CacheRecords(records []models.Record) error {
	data, err := json.Marshal(records)
	if err != nil {
		return err
	}
	return config.RedisClient.Set(context.Background(), "imported_data", data, 5*time.Minute).Err()
}

// GetRecords fetches the cached records from Redis
func GetRecords() ([]models.Record, string, error) {
	var records []models.Record
	cached, err := config.RedisClient.Get(context.Background(), "imported_data").Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cached), &records); err == nil {
			return records, "cache", nil
		}
	}
	return nil, "", errors.New("no data found")
}

// UpdateRecord updates a record in the Redis cache
func UpdateRecord(email string, updated models.Record) error {
	records, _, err := GetRecords()
	if err != nil {
		return err
	}
	for i, rec := range records {
		if rec.Email == email {
			records[i] = updated
			break
		}
	}
	return CacheRecords(records)
}
