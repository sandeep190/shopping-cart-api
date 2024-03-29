package models

import (
	"html/template"

	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

type Products struct {
	gorm.Model
	Title    string `gorm:"size:255;not null"`
	Details  string `gorm:"not null"`
	SortDesc string `gorm:"not null" text:"html"`
	Slug     string `gorm:"size:255;unique_index;not null"`
	Price    uint32 `gorm:"not null"`
	Quantity int    `gorm:"not null"`
	CatID    int    `gorm:"not null"`
	SubcatID int    `gorm:"not null"`
	// Categories []Category `gorm:"many2many:products_categories;"`

	Images string `gorm:"not null"`
}

func (product *Products) BeforeSave(*gorm.DB) (err error) {
	product.Slug = slug.Make(product.Title)
	return
}

type ProductList struct {
	ID       int
	Title    string
	SortDesc template.HTML
	CatID    int
	SubcatID int
	Details  template.HTML
	Price    float64
	Quantity int
	Images   string
}
