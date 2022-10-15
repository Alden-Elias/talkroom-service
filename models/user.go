package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Name        string  `gorm:"type:varchar(35);not null;"`
	Avatar      *string `gorm:"type:mediumtext;"`
	Email       string  `gorm:"unique;type:varchar(35);not null;"`
	Password    string  `gorm:"type:varchar(120);not null;"`
	Description string  `gorm:"type:varchar(300)"`
	Sex         string  `gorm:"type:char(3)"`
	Birthday    time.Time
	IsAdmin     bool `gorm:"default:false;not null;"`
	IsLocked    bool `gorm:"default:false;not null;"`
}

type UserBrief struct {
	ID     uint    `json:"id"`
	Name   string  `json:"name"`
	Email  string  `json:"email"`
	Avatar *string `json:"avatar"`
}

type UserDetailed struct {
	UserBrief
	Description string `json:"description"`
	Sex         string `json:"sex"`
	Birthday    string `json:"birthday"`
}
