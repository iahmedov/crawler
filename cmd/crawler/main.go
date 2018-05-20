package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/iahmedov/crawler"
	"gopkg.in/yaml.v2"

	_ "github.com/iahmedov/crawler/fetcher/circuitbreaker"
	_ "github.com/iahmedov/crawler/fetcher/limit/body"
	_ "github.com/iahmedov/crawler/filter/content/contenttype"
	_ "github.com/iahmedov/crawler/filter/content/porn"
	_ "github.com/iahmedov/crawler/filter/link/depth"
	_ "github.com/iahmedov/crawler/filter/link/domain"
	_ "github.com/iahmedov/crawler/filter/link/social/facebook"
	_ "github.com/iahmedov/crawler/filter/link/storage"
	_ "github.com/iahmedov/crawler/link/html"
	_ "github.com/iahmedov/crawler/queue/local"
	_ "github.com/iahmedov/crawler/storage/memory"
	_ "github.com/iahmedov/crawler/task/transition/unreliable"
)

func main() {
	file, err := os.Open("config.yaml")
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	buff, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err.Error())
	}

	config := crawler.Config{}
	if err = yaml.Unmarshal(buff, &config); err != nil {
		panic(err.Error())
	}

	validator := crawler.ValidateConfig(config)
	if validator != nil && validator.HasError() {
		fmt.Printf("errors: \n%s\n", validator.Error())
		return
	}

	cr, err := crawler.NewCrawler(config)
	if err != nil {
		panic(err.Error())
	}

	if err := cr.Run(context.Background()); err != nil {
		panic(err.Error())
	}
}
