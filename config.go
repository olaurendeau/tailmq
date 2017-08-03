package main

import (
  "io/ioutil"
  "errors"
  
  "github.com/mitchellh/go-homedir"
  "gopkg.in/yaml.v2"
)

type GlobalConfig struct {
  Servers map[string]string
}

func getServerList(path string) (GlobalConfig, error) {

  var globalConfig GlobalConfig

  if path == "" {
    homedir, err := homedir.Dir()
    if err != nil {
      return globalConfig, err
    }

    path = homedir + "/.tailmq"
  }

  fileContent, err := ioutil.ReadFile(path)

  if err != nil {
    return globalConfig, err
  }
  
  err = yaml.Unmarshal(fileContent, &globalConfig)

  return globalConfig, err
}

func getServerUri(server string, globalConfig GlobalConfig) (string, error) {
  var err error
  uri := globalConfig.Servers[server]
  if uri == "" {
    err = errors.New("No server named " + server)
  }

  return uri, err
}