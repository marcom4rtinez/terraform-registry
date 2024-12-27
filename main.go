package main

import (
	"github.com/gin-gonic/gin"
	routes "github.com/marcom4rtinez/terraform-registry/router"
)

func main() {
	r := gin.Default()

	routes.RegisterRoutes(r)

	r.Run()
}
