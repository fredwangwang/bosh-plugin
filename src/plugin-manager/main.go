package main

import (
	"github.com/cloudfoundry/go-envstruct"
	"github.com/fredwangwang/bosh-plugin/pluginmanager"
	"github.com/fredwangwang/bosh-plugin/routes"
	"github.com/gin-gonic/gin"
	"log"
)

type HostInfo struct {
	Port         string `env:"PORT, required"`
	Job          string `env:"JOB, required"`
	Monit        string `env:"MONIT, required"`
	Storage      string `env:"STORAGE, required"`
	PluginConfig string `env:"PLUGIN_CONFIG_FILE, required"`

	ServerCertPath string `env:"SERVER_CERT_PATH, required"`
	ServerKeyPath  string `env:"SERVER_KEY_PATH, required"`

	UAAUrl string   `env:"UAA_URL, required"`
	Scopes []string `env:"ALLOWED_SCOPES, required"`
}

func main() {
	info := HostInfo{}
	if err := envstruct.Load(&info); err != nil {
		panic(err)
	}

	pm, err := pluginmanager.GetPluginManager(info.Job, info.Monit, info.Storage, info.PluginConfig)
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	routes.RegistRoutes(r, pm, info.UAAUrl, info.Scopes)

	done := make(chan bool)

	go func() {
		r.RunTLS(":"+info.Port, info.ServerCertPath, info.ServerKeyPath)
		done <- true
	}()

	log.Println("plugin server is running!")

	<-done
}
