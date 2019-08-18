package routes

import (
	"fmt"
	"github.com/fredwangwang/bosh-plugin/pluginmanager"
	"github.com/fredwangwang/bosh-plugin/uaa"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"path"
)

func RegistRoutes(r *gin.Engine, pm pluginmanager.Manager, uaaUrl string, scopes []string) {
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "I am up!",
		})
	})

	authed := r.Group("/api")
	authed.Use(uaa.UAAJWTMiddleware(uaaUrl, scopes))
	{
		authed.GET("/plugins", handleList(pm))

		authed.GET("/plugins/:name", handleGet(pm))

		authed.POST("/plugins", handleUpload(pm))

		authed.POST("/plugins/:name/enable", handleEnable(pm))

		authed.POST("/plugins/:name/disable", handleDisable(pm))

		authed.PATCH("/plugins/:name", handleConfig(pm))

		authed.DELETE("/plugins/:name", handleDelete(pm))
	}
}

func handleList(pm pluginmanager.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		if stats, err := pm.ListPlugins(); err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, stats)
		}
	}
}

func handleGet(pm pluginmanager.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		pluginName := c.Param("name")
		if state, err := pm.GetPlugin(pluginName); err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, state)
		}
	}
}

func handleUpload(pm pluginmanager.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, _ := c.FormFile("file")
		log.Println(file.Filename + " received")

		log.Println("creating tmp dir")
		tempdir, err := ioutil.TempDir("", "uploaded")
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		filename := path.Join(tempdir + "uploaded.zip")
		log.Println("save uploaded file to " + filename)
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := pm.AddPlugin(filename); err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "plugin uploaded successfully",
		})
	}
}

func handleEnable(pm pluginmanager.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		pluginName := c.Param("name")
		if err := pm.EnablePlugin(pluginName); err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, fmt.Sprintf("%s enabled successfully", pluginName))
		}
	}
}

func handleDisable(pm pluginmanager.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		pluginName := c.Param("name")
		if err := pm.DisablePlugin(pluginName); err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, fmt.Sprintf("%s disabled successfully", pluginName))
		}
	}
}

func handleConfig(pm pluginmanager.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		pluginName := c.Param("name")
		if err := pm.ConfigPlugin(pluginName, c.Request.URL.Query()); err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, fmt.Sprintf("config applied, need to disable/enable the plugin to see the effect"))
		}
	}
}

func handleDelete(pm pluginmanager.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		pluginName := c.Param("name")
		if err := pm.DeletePlugin(pluginName); err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, fmt.Sprintf("%s deleted successfully", pluginName))
		}
	}
}
