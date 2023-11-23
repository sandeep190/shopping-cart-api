package controllers

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Pages struct {
	Title   string        `title`
	Content template.HTML `content`
}

func PagesRoutes(routes *gin.RouterGroup) {
	routes.GET("", Index)
	routes.GET("/pages/about", AboutUs)
}

func AboutUs(c *gin.Context) {
	data := make(map[string]string)
	data["title"] = "About Us page "
	data["content"] = "This Is the About Us Page"
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
