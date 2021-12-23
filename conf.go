package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Configuration struct {
	Name        string `yaml:"name"`
	LocalIp     string `yaml:"localIp"`
	ChannelPort int    `yaml:"channelPort"`
	Socks5Port  int    `yaml:"socks5Port"`
	Gateway     string `yaml:"gateway"`
	Ping        int    `yaml:"ping"`
}

func GetConf(path string) Configuration {
	var conf Configuration
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(content, &conf)
	if err != nil {
		panic(err)
	}
	return conf
}
