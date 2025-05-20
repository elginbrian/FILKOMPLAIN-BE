package model

import "gorm.io/gorm"

type Report struct {
	gorm.Model
	UserName    string `json:"user_name"`
	Content     string `json:"content"`
	Place       string `json:"place"`
	PhoneNumber string `json:"phone_number"`
	Status      string `json:"status"`
	Attachment  string `json:"attachment"`
	Reply       string `json:"reply"` 
}
