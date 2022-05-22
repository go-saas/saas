package main

import (
	gorm2 "github.com/goxiaoy/go-saas/gorm"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	//Do not use in real case!
	DSN string `json:"DSN"`
	gorm2.MultiTenancy
}
