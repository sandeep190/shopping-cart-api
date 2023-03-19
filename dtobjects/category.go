package dtobjects

import (
	"shopping_cart/models"
	"strings"
)

type SaveCategoryRequestDto struct {
}

func CategoryListDto(categories []models.Category) map[string]interface{} {
	result := map[string]interface{}{}
	var t = make([]interface{}, len(categories))
	for i := 0; i < len(categories); i++ {
		t[i] = CreateCategoryDto(categories[i])
	}
	result["categories"] = t
	return CreateSuccessDto(result)
}

func CreateCategoryDto(category models.Category) map[string]interface{} {
	var imageUrls = make([]string, len(category.Images))
	replaceAllFlag := -1
	for i := 0; i < len(category.Images); i++ {
		imageUrls[i] = strings.Replace(category.Images[i].FilePath, "\\", "/", replaceAllFlag)
	}
	return map[string]interface{}{
		"id":          category.ID,
		"name":        category.Name,
		"parent_id":   category.ParentId,
		"description": category.Description,
		"image_urls":  imageUrls,
	}
}

func CategoryListAdminDto(categories []models.CatagoryList) map[string]interface{} {
	result := map[string]interface{}{}
	var t = make([]interface{}, len(categories))
	for i := 0; i < len(categories); i++ {
		t[i] = CreateCategoryAdminDto(categories[i])
	}
	result["categories"] = t
	return CreateSuccessDto(result)
}

func CreateCategoryAdminDto(category models.CatagoryList) map[string]interface{} {
	imageUrls := strings.Replace(category.FilePath, "\\", "/", -1)
	return map[string]interface{}{
		"id":          category.ID,
		"name":        category.Name,
		"parent_id":   category.ParentId,
		"description": category.Description,
		"image_urls":  imageUrls,
		"parent":      category.Parent,
	}
}
