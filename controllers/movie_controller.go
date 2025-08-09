package controllers

import (
	"backend/confiq"
	"backend/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateMovie(c *gin.Context) {
	title := c.PostForm("title")
	description := c.PostForm("description")

	if title == "" || description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title and description are required"})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		log.Printf("Image file is required: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}

	if file.Size > 2<<20 { 
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 2MB"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only .jpg or .jpeg files are allowed"})
		return
	}

	
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	imagePath := filepath.Join("uploads", fileName)

	if err := c.SaveUploadedFile(file, imagePath); err != nil {
		log.Printf("Failed to save uploaded file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
		return
	}

	
	movie := models.Movie{
		Title:       title,
		Description: description,
		Image:       fileName, 
	}

	if err := confiq.DB.Create(&movie).Error; err != nil {
		log.Printf("Failed to insert movie: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	
	imageURL := fmt.Sprintf("http://%s/uploads/%s", c.Request.Host, movie.Image)

	log.Printf("Movie created successfully: %+v", movie)
	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"id":          movie.ID,
			"title":       movie.Title,
			"description": movie.Description,
			"image":       imageURL,
		},
	})
}



func GetMovie(c *gin.Context) {
	var movies []models.Movie
	if err := confiq.DB.Find(&movies).Error; err != nil {
		log.Printf("terjadi error:%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "terjadi error"})
		return
	}
	if len(movies) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "data not found"})
		return
	}

	var response []map[string]interface{}
	for _, m := range movies {
		imageURL := fmt.Sprintf("http://%s/uploads/%s", c.Request.Host, m.Image)
		response = append(response, map[string]interface{}{
			"id":          m.ID,
			"title":       m.Title,
			"description": m.Description,
			"image":       imageURL,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}



func GetByID(c *gin.Context) {
	id := c.Param("id")

	var movie models.Movie
	if err := confiq.DB.First(&movie, id).Error; err != nil {
		log.Printf("Movie with ID %s not found: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
		return
	}

	imageURL := fmt.Sprintf("http://%s/uploads/%s", c.Request.Host, movie.Image)

	response := map[string]interface{}{
		"id":          movie.ID,
		"title":       movie.Title,
		"description": movie.Description,
		"image":       imageURL,
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}


func UpdateMovie(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("Invalid ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var movie models.Movie
	if err := confiq.DB.First(&movie, uint(id)).Error; err != nil {
		log.Printf("Movie not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")


	movie.Title = title
	movie.Description = description

	
	file, err := c.FormFile("image")
	if err == nil {
		
		oldPath := filepath.Join("uploads", movie.Image)
		if err := os.Remove(oldPath); err != nil {
			log.Printf("Failed to delete old image: %v", err)
		}

		
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only .jpg or .jpeg files are allowed"})
			return
		}

		
		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		savePath := filepath.Join("uploads", newFileName)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			log.Printf("Failed to upload new image: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload new image"})
			return
		}

		movie.Image = newFileName
	}


	if err := confiq.DB.Save(&movie).Error; err != nil {
		log.Printf("Failed to update movie: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie"})
		return
	}

	
	imageURL := fmt.Sprintf("http://%s/uploads/%s", c.Request.Host, movie.Image)
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":          movie.ID,
			"title":       movie.Title,
			"description": movie.Description,
			"image":       imageURL,
		},
	})
}

func DeleteMovie(c *gin.Context) {
	id := c.Param("id")

	var movie models.Movie
	if err := confiq.DB.First(&movie, id).Error; err != nil {
		log.Printf("Movie with ID %s not found: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
		return
	}

	// Hapus file gambar dari folder uploads/
	imagePath := filepath.Join("uploads", movie.Image)
	if err := os.Remove(imagePath); err != nil && !os.IsNotExist(err) {
		log.Printf("Failed to delete image file: %v", err)
		// lanjutkan penghapusan dari DB meskipun file gagal dihapus
	}

	// Hapus data dari DB
	if err := confiq.DB.Delete(&movie).Error; err != nil {
		log.Printf("Failed to delete movie from DB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie"})
		return
	}

	log.Printf("Movie with ID %s deleted successfully", id)
	c.JSON(http.StatusOK, gin.H{"message": "Movie deleted successfully"})
}
