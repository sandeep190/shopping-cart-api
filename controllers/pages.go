package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PagesRoutes(routes *gin.RouterGroup) {
	routes.GET("", Index)
	routes.GET("/pages/about", AboutUs)
}

func AboutUs(c *gin.Context) {
	data := make(map[string]string)
	data["title"] = "index page "
	data["content"] = "this is the index page"
	c.HTML(http.StatusOK, "about.html", gin.H{
		"content": data,
	})
}

func Index(c *gin.Context) {

	data := make(map[string]string)
	data["title"] = "index page "
	data["content"] = "this is the index page"
	c.HTML(http.StatusOK, "index.html", gin.H{
		"content": data,
	})
}
