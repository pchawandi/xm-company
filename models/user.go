package models

import "time"

type LoginUser struct {
	Username string `json:"username" binding:"required"`
	Role     string `json:"role" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Username  string    `json:"username" gorm:"unique"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
