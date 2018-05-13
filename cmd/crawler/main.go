package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/iahmedov/crawler"
	"gopkg.in/yaml.v2"
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
		fmt.Println("errors: %s", validator.Error())
	}
}
