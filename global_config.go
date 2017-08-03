package main

import (
  "io/ioutil"
  "errors"
  "os/user"
  
  "gopkg.in/yaml.v2"
)

type GlobalConfig struct {
  Servers map[string]string
}

func NewGlobalConfig(path string) (GlobalConfig, error) {

  var globalConfig GlobalConfig

  if path == "" {
    usr, err := user.Current()
    if err != nil {
      return globalConfig, err
    }

    path = usr.HomeDir + "/.tailmq"
  }

  fileContent, err := ioutil.ReadFile(path)

  if err != nil {
    return globalConfig, err
  }
  
  err = yaml.Unmarshal(fileContent, &globalConfig)

  return globalConfig, err
}

func (g GlobalConfig) getServerUri(server string) (string, error) {
  var err error
  uri := g.Servers[server]
  if uri == "" {
    err = errors.New("No server named " + server)
  }

  return uri, err
}