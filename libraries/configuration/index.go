package configuration

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"encoding/json"
)

type ConfigType struct {
	IpAddress	string	`json:"ip_address"`
	Port			string	`json:"port"`
}

func (config *ConfigType) Save() {
	data, err := json.Marshal(config);
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(filepath.FromSlash(fmt.Sprintf("%s%s%s", RootPath, "/", "config.json")), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	if _, err := f.Write(data); err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

var (
	RootPath string
	Config ConfigType = ConfigType{IpAddress: "",Port: ""}
)

func init() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	RootPath = filepath.Dir(ex)
	configPath := filepath.FromSlash(fmt.Sprintf("%s%s%s", RootPath, "/", "config.json"))
	log.Print(configPath)
	if data, err := os.ReadFile(configPath); err != nil {
		log.Print(err)
		json := []byte("{}")
		if err := os.WriteFile(configPath, json, 0644); err != nil {
			panic(err)
		} else {
			if _, err := os.ReadFile(configPath); err != nil {
				panic(err)
			}
		}
	} else {
		if err := json.Unmarshal(data, &Config); err != nil {
			panic(err)
		}
	}
	log.Print(Config.IpAddress, ":", Config.Port)
}
