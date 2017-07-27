package main

import (
  "fmt"
  "errors"
  "io/ioutil"
  "encoding/json"

  "github.com/mitchellh/go-homedir"
)

type AmqpConfig struct {
    Name string`json:"name"`
    Uri string `json:"uri"`
}

func getAmqpConfigList() []AmqpConfig {
  homedir, err := homedir.Dir()
  if err != nil {
      fmt.Print(err)
  }

  fileContent, err := ioutil.ReadFile(homedir + "/.tailmq")
  
  if err != nil {
      fmt.Print(err)
  }

  var configList []AmqpConfig
  json.Unmarshal(fileContent, &configList)

  return configList
}

func getConfig(name string, configList []AmqpConfig) (AmqpConfig, error) {

  var config AmqpConfig

  for _, config := range configList {
    if (config.Name == name) {
      return config, nil
    }
  }

  return config, errors.New("No config named " + name)
}
