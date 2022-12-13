package models

import "github.com/jinzhu/gorm"

type FileUpload struct {
	gorm.Model
	Filename     string
	FilePath     string
	OriginalName string
	FileSize     uint

	Category   Category `gorm:"association_foreignkey:CategoryId"`
	CategoryId uint     `gorm:"default:null"`

	Product   Category `gorm:"association_foreignkey:ProductId"`
	ProductId uint     `gorm:"default:null"`
}

func CategoryImages(db *gorm.DB) *gorm.DB {
	return db.Where("type = ?", "CategoryImage")
}

func ProductImages(db *gorm.DB) *gorm.DB {
	return db.Where("type = ?", "ProductImage")
}
