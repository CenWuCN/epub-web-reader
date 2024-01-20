package setting

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	EpubsPath          string `yaml:"epubsPath"`
	UnzipPath          string `yaml:"unzipPath"`
	GinEpubsStaticPath string `yaml:"ginEpubsStaticPath"`
}

var EpubsAbsPath string = ""
var UnzipAbsPath string = ""
var ConfigYaml = Config{}

func Init() {
	bytes, err := os.ReadFile("./config.yaml")
	if err != nil {
		fmt.Println(err)
	}
	yaml.Unmarshal(bytes, &ConfigYaml)

	EpubsAbsPath, err = ExpandTilde(ConfigYaml.EpubsPath)
	if err != nil {
		fmt.Println("ExpandTilde Err", ConfigYaml.EpubsPath, err)
	}
	UnzipAbsPath, err = ExpandTilde(ConfigYaml.UnzipPath)
	if err != nil {
		fmt.Println("ExpandTilde Err", ConfigYaml.UnzipPath, err)
	}
}

func ExpandTilde(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Join(usr.HomeDir, path[1:]), nil
}
