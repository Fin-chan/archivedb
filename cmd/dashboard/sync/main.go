package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/yaml.v3"

	"github.com/jialeicui/archivedb/pkg"
)

type Config struct {
	Uid    string `yaml:"uid"`
	Cookie string `yaml:"cookie"`
}

const (
	configPath = ".config.yaml"
	dbPath     = ".data"
)

func main() {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	config := new(Config)
	err = yaml.Unmarshal(content, config)
	if err != nil {
		panic(err)
	}

	os.RemoveAll(dbPath)
	db, err := pkg.New(dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	cli, err := NewWithHeader(map[string]string{"cookie": config.Cookie})
	if err != nil {
		panic(err)
	}

	url := fmt.Sprintf("https://weibo.com/ajax/favorites/all_fav?uid=%s&page=2", config.Uid)
	resp, err := cli.Get(url)
	if err != nil {
		panic(err)
	}

	item := pkg.Item{}
	err = bson.UnmarshalExtJSON(resp, true, &item)
	if err != nil {
		panic(err)
	}
	ok := item["ok"]
	if ok != int32(1) {
		panic(fmt.Sprintf("invalid content: %q", string(resp)))
	}
	for _, i := range item["data"].(bson.A) {
		it := i.(bson.M)
		err = db.Put(&it,
			pkg.WithKey([]byte(it["idstr"].(string))),
		)
		if err != nil {
			panic(err)
		}
	}
}
