package main

import (
	"check-links/models"
	"check-links/storage"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var store *storage.JSONStorage

func main() {
	var err error
	store, err = storage.NewJSONStorage("data.json")
	if err != nil {
		log.Fatal(err)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	public := r.Group("api/v1")
	{
		public.POST("/links/check", LinksCheck)
		public.GET("/links", StatusList)
		public.PUT("/linklists/:num/status", updateLinkListStatus)
	}
}

func LinksCheck(c *gin.Context) {
	var input struct {
		Links []string `json:"links"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	linkObj := make([]models.Link, len(input.Links))
	for i, link := range input.Links {
		status := models.LinkStatusUnknown

		if checkLinkStatus(link) {
			status = models.LinkStatusAvailable
		} else {
			status = models.LinkStatusNotAvailable
		}

		linkObj[i] = models.Link{
			Link:   link,
			Status: status,
		}
	}

	err := store.Create(models.LinkList{Links: linkObj})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to save links"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"links": linkObj})
}

func StatusList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"links": store.List()})
}

func checkLinkStatus(link string) bool {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(link)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 400
}

func updateLinkListStatus(c *gin.Context) {
	numStr := c.Param("num")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid num"})
		return
	}

	if err := store.UpdateStatus(num); err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "list not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}
