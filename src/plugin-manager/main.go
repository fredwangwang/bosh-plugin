package main

import (
	"github.com/cloudfoundry/go-envstruct"
	"github.com/fredwangwang/bosh-plugin-manager/pluginmanager"
	"github.com/fredwangwang/bosh-plugin-manager/routes"
	"github.com/gin-gonic/gin"
	"log"
)

type HostInfo struct {
	Port         string `env:"PORT, required"`
	Storage      string `env:"STORAGE, required"`
	PluginConfig string `env:"PLUGIN_CONFIG_FILE, required"`
}

func main() {
	info := HostInfo{}
	if err := envstruct.Load(&info); err != nil {
		panic(err)
	}

	pm := pluginmanager.GetPluginManager(info.Storage, info.PluginConfig)

	r := gin.Default()
	routes.RegistRoutes(r, pm)

	done := make(chan bool)

	go func() {
		r.Run(":" + info.Port)
		done <- true
	}()

	log.Println("plugin server is running!")

	<-done
}
