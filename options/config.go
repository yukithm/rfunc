package options

import (
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml"
	"github.com/spf13/pflag"
	"github.com/yukithm/rfunc/utils"
)

var defaultConfigFiles = []string{
	"~/.config/rfunc/rfunc.toml",
	"~/.rfunc.toml",
}
var ConfigFile string
var configFlagSet *pflag.FlagSet

func init() {
	configFlagSet = pflag.NewFlagSet(filepath.Base(os.Args[0]), pflag.ContinueOnError)
	configFlagSet.ParseErrorsWhitelist.UnknownFlags = true
	configFlagSet.StringVarP(&ConfigFile, "conf", "c", ConfigFile, "configuration file")
}

func LoadConfig() (*Options, error) {
	configFlagSet.Parse(os.Args[1:])
	if ConfigFile == "" {
		ConfigFile = findDefaultConfigFile()
	}

	if ConfigFile != "" {
		return loadConfigFile(ConfigFile)
	}

	return &Options{}, nil
}

func loadConfigFile(conf string) (*Options, error) {
	path, err := utils.ExpandPath(conf)
	if err != nil {
		return nil, err
	}
	return loadTOMLConfig(path)
}

func findDefaultConfigFile() string {
	for _, file := range defaultConfigFiles {
		path, err := utils.ExpandPath(file)
		if err != nil {
			return ""
		}

		if utils.FileExists(path) {
			return path
		}
	}

	return ""
}

func loadTOMLConfig(file string) (*Options, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var conf Options
	decoder := toml.NewDecoder(f)
	if err := decoder.Decode(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
