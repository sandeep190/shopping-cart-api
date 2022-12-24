package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(routes *gin.RouterGroup) {
	routes.GET("admin", AdminIndex)
	routes.GET("/admin/master/category", AdminCategoryList)
}

func AdminIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_index.html", gin.H{
		"content": "Admin page",
	})
}

func AdminCategoryList(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_category.html", gin.H{
		"content": "Admin page",
	})
}
