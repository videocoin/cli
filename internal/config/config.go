package config

import (
	"fmt"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type Config struct {
	ManagerAddr     string `required:"true" default:""`
	NodeRPCAddr     string `required:"true" envconfig:"NodeRPCAddr" default:""`
	ContractAddress string `required:"true" envconfig:"ContractAddress" default:"" description:"stream manager contract address"`

	MinVDCBalance int `required:"true" default:"15"`

	Logger   *logrus.Entry `ignored:"true"`
	Loglevel string        `default:"FATAL" envconfig:"LOGLEVEL"`
}

func (c *Config) InitLogger(name, version string) error {
	level, err := logrus.ParseLevel(c.Loglevel)
	if err != nil {
		return fmt.Errorf("not a valid log level: %q", c.Loglevel)
	}

	logrus.SetLevel(level)
	//logrus.AddHook(filename.NewHook())

	logger := logrus.WithFields(logrus.Fields{
		"service": name,
		"version": version,
	})

	formatter := new(prefixed.TextFormatter)
	formatter.SetColorScheme(&prefixed.ColorScheme{
		PrefixStyle:    "green+b",
		TimestampStyle: "white+h",
	})

	logger.Logger.Formatter = formatter
	c.Logger = logger

	return nil
}
