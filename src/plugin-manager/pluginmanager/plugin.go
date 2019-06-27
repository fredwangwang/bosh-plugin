package pluginmanager

import (
	"archive/zip"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

const PLUGIN_INFO_FILE = "plugin.yml"

type Manager struct {
	Storage    string
	ConfigFile string
}

type State struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Location    string `yaml:"location"`
	Enabled     bool   `yaml:"enabled"`
}

type States []State

func GetPluginManager(storage string, configFile string) Manager {
	pm := Manager{
		Storage:    storage,
		ConfigFile: configFile,
	}

	return pm.init()
}

func (p Manager) ListPlugins() (States, error) {
	var states States
	var err error

	if configBytes, err := ioutil.ReadFile(p.configFilePath()); err != nil {
		return states, errors.Wrap(err, "failed to read config file")
	} else {
		err = yaml.Unmarshal(configBytes, &states)
	}

	return states, errors.Wrap(err, "failed to unmarshal config file")
}

func (p Manager) AddPlugin(filename string) error {
	var err error

	// TODO: handle dup with existing ones
	states, err := p.ListPlugins()
	if err != nil {
		return errors.Wrap(err, "failed to retrieve existing plugins")
	}

	infos, err := getPluginInfo(filename)
	if err != nil {
		return err
	}

	info := infos.Applications[0]

	location := strings.Join(strings.Fields(info.Name), "")
	pluginPath := path.Join(p.Storage, location)

	cmd := exec.Command("unzip", filename, "-d", pluginPath)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed to unzip plugin")
	}

	// TODO: populate executor (BPM) context

	states = append(states, State{
		Name:        info.Name,
		Description: info.Description,
		Location:    location,
		Enabled:     false,
	})

	configFH, err := os.Open(p.configFilePath())
	if err != nil {
		return err
	}
	defer configFH.Close()

	return yaml.NewEncoder(configFH).Encode(states)
}

func (p Manager) DeletePlugin(pluginName string) error {
	if _, err := os.Stat(p.pluginPath(pluginName)); os.IsNotExist(err) {
		return fmt.Errorf("plugin %s does not exist", pluginName)
	}

	// TODO: stop the plugin
	// TODO: cleanup executor context

	return errors.Wrap(os.RemoveAll(p.pluginPath(pluginName)),
		"failed to remove plugin content")
}

func (p Manager) configFilePath() string {
	configFilePath := path.Join(p.Storage, p.ConfigFile)
	return configFilePath
}

func (p Manager) pluginPath(pluginName string) string {
	return path.Join(p.Storage, pluginName)
}

func (p Manager) init() Manager {
	if _, err := os.Stat(p.configFilePath()); os.IsNotExist(err) {
		log.Println("initializing " + p.ConfigFile)
	}

	fh, err := os.Create(p.configFilePath())
	if err != nil {
		panic(err)
	}
	fh.Close()

	return p
}

type Application struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Command     string            `yaml:"command"`
	Memory      string            `yaml:"memory"`
	Args        []string          `yaml:"args"`
	Env         map[string]string `yaml:"env"`
}

type Info struct {
	Applications []Application `yaml:"applications"`
}

func getPluginInfo(s string) (Info, error) {
	info := Info{}

	r, err := zip.OpenReader(s)
	if err != nil {
		return info, errors.Wrap(err, "failed to open plugin file")
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name != PLUGIN_INFO_FILE {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return info, errors.Wrap(err, "failed to read plugin info")
		}
		err = yaml.NewDecoder(rc).Decode(&info)
		if err != nil {
			return info, errors.Wrap(err, "failed to unmarshal plugin info")
		}

		rc.Close()

		// TODO: verify only one application entry is provided
		return info, ValidatePluginName(info.Applications[0].Name)
	}
	return info, fmt.Errorf(PLUGIN_INFO_FILE + " does not exist in the plugin file")
}

var IsValidName = regexp.MustCompile(`^[a-zA-Z0-9_\- ]+$`).MatchString

func ValidatePluginName(name string) error {
	if !IsValidName(name) {
		return fmt.Errorf("%s is not a valid plugin name, only 'a-z A-Z 0-9 _ -' are allowed", name)
	}
	return nil
}
