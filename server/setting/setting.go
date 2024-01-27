package setting

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	EpubsPath          string   `yaml:"epubsPath"`
	UnzipPath          string   `yaml:"unzipPath"`
	GinEpubsStaticPath string   `yaml:"ginEpubsStaticPath"`
	Jwtkey             string   `yaml:"jwtkey"`
	Invitecodes        []string `yaml:"invitecode"`
}

var EpubsAbsPath string = ""
var UnzipAbsPath string = ""
var JwtkeyBytes []byte
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
	if ConfigYaml.Jwtkey == "" {
		bytes = make([]byte, 32)
		_, err = rand.Read(bytes)
		if err != nil {
			fmt.Println(err)
		}
		jwtkey := base64.StdEncoding.EncodeToString(bytes)
		ConfigYaml.Jwtkey = jwtkey
		bytes, err = yaml.Marshal(&ConfigYaml)
		if err != nil {
			fmt.Println(err)
		}
		os.WriteFile("./config.yaml", bytes, os.ModePerm)
	}
	JwtkeyBytes = []byte(ConfigYaml.Jwtkey)
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
