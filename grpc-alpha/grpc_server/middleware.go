package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TempURLHandler Checks the DB to verify that the url results in a file
func TempURLHandler(c *gin.Context) {
	fileID := c.Param("file_id")
	fileInfo := TempURL{}
	dbResult := db.First(&fileInfo, map[string]interface{}{"url_post": fileID, "valid": true})
	if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatus(404)
		return
	}
	c.File(fileInfo.FilePath)
	fileInfo.Valid = false
	db.Save(&fileInfo)
	return
}

// TaskAdder Adds task to DB
func TaskAdder(c *gin.Context) {
	var json Task
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dbResult := db.Create(&json)
	if dbResult.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Task has failed to be created"})
		return

	}
	c.JSON(http.StatusOK, gin.H{"status": "Task Created"})
	return
}

// FileAdder Adds File Path to DB
func FileAdder(c *gin.Context) {
	var json TempURL
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	json.URLPost = uuid.New().String()
	dbResult := db.Create(&json)
	if dbResult.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "File URL has failed to be created"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "File URL was Created"})
	return
}
