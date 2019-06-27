package routes

import (
	"github.com/fredwangwang/bosh-plugin-manager/pluginmanager"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
)

func RegistRoutes(r *gin.Engine, pm pluginmanager.Manager) {
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "I am up!",
		})
	})

	r.GET("/plugins", func(c *gin.Context) {
		if stats, err := pm.ListPlugins(); err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, stats)
		}
	})

	r.POST("/plugins", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		log.Println(file.Filename + " received")
		tempdir, err := ioutil.TempDir("uploaded", "")
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		filename := tempdir + "uploaded.zip"

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
	})

}

// /var/vcap/store
