package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Uid    string `yaml:"uid"`
	Cookie string `yaml:"cookie"`
}

func main() {
	configPath := ".config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	config := new(Config)
	err = yaml.Unmarshal(content, config)
	if err != nil {
		panic(err)
	}

	cli, err := NewWithHeader(map[string]string{"cookie": config.Cookie})
	if err != nil {
		panic(err)
	}

	url := fmt.Sprintf("https://weibo.com/ajax/favorites/all_fav?uid=%s&page=1", config.Uid)
	resp, err := cli.Get(url)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(resp))
}
