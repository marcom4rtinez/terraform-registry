package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marcom4rtinez/terraform-registry/controller"
)

func RegisterRoutes(router *gin.Engine) {

	provider := router.Group("/v1/providers")
	{
		provider.GET("/:namespace/:name/:version/download/:os/:arch", controller.DownloadProvider)
		provider.GET("/:namespace/:name/versions", controller.GetVersion)
		provider.POST("/:namespace/:name/upload", controller.UploadProvider)
	}

	router.GET("/.well-known/terraform.json", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"providers.v1": "/v1/providers/"})
	})
}
