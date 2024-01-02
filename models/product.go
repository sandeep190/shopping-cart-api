package models

import (
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

type Product struct {
	gorm.Model
	Title      string  `gorm:"size:255;not null"`
	Detail     string  `gorm:"not null"`
	SortDesc   string  `gorm:"not null"`
	Slug       string  `gorm:"size:255;unique_index;not null"`
	Price      float32 `gorm:"not null"`
	Quantity   int     `gorm:"not null"`
	CategoryID int     `gorm:"not null"`

	Categories []Category `gorm:"many2many:products_categories;"`

	Images string `gorm:"not null"`
}

func (product *Product) BeforeSave(*gorm.DB) (err error) {
	product.Slug = slug.Make(product.Title)
	return
}

type ProductList struct {
	ID       int
	Title    string
	SortDesc string
	CatID    int
	Details  string
	Price    float32
	Quantity int
	image    string
}
