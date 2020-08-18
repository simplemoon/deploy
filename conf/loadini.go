package conf

import (
	"gopkg.in/ini.v1"
)

func LoadIni(name string) (*ini.File, error) {
	cfg, err := ini.Load(name)
	if err != nil {
		return nil, err
	}
	return cfg, err
}
