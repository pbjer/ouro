package cli

import "gorm.io/gorm"

type FileContent struct {
	gorm.Model
	Path    string
	Content string `gorm:"type:text;"`
}
