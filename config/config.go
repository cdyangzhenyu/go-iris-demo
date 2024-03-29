package config

import (
	"fmt"
	"github.com/pelletier/go-toml"
)

var (
	Conf = New()
)

/**
 * 返回单例实例
 * @method New
 */
func New() *toml.Tree {
	config, err := toml.LoadFile("/etc/bms/config.ini")

	if err != nil {
		fmt.Println("TomlError ", err.Error())
	}

	return config
}
