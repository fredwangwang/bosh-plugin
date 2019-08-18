package main

import (
	"github.com/gin-gonic/gin"
	"test-gin/uaa-middleware"
)

func main() {
	// Creates a router without any middleware by default
	r := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	authorized := r.Group("/auth")

	authorized.Use(uaa_middleware.UAAJWTMiddleware(
		"https://login.moonstorm.cf-denver.com",
		[]string{"cloud_controller.admin"},
	))

	authorized.GET("/test", func(c *gin.Context) {
		c.JSON(200, "hello!")
	})
	// Listen and serve on 0.0.0.0:8080
	r.Run("localhost:8080")
}
