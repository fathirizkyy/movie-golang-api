package models

type Movie struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Title       string `json:"title"`
	Image       string `json:"image"` // nama file .jpg
	Description string `json:"description"`
}