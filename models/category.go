package models

import (
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

type Category struct {
	gorm.Model
	ID          int          `gorm:"primary key"`
	Name        string       `gorm:"not null"`
	Description string       `gorm:"default:null"`
	ParentId    int          `gorm:"default:0"`
	Image       string       `gorm:"default:null"`
	Slug        string       `gorm:"unique_index"`
	Products    []Product    `gorm:"many2many:products_categories;"`
	Images      []FileUpload `gorm:"foreignKey:CategoryId"`
	IsNewRecord bool         `gorm:"-;default:false"`
}

func (cat *Category) BeforeSave(*gorm.DB) (err error) {
	cat.Slug = slug.Make(cat.Name)
	return
}
