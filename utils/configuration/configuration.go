package configuration

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"encoding/json"
	"net"
	"strings"
	"errors"
	"strconv"
)

type ConfigType struct {
	IpAddress	string	`json:"ip_address"`
	Port			string	`json:"port"`
}

func CheckIPAddress(ip, port string) (string, error) {
	ipAddr := strings.Join([]string{ip, port}, ":")
	if net.ParseIP(ip) == nil {
		return ipAddr, errors.New(strings.Join([]string{"Error:", ip, "is invalid IP address"}, " "))
	}
	p, err := strconv.Atoi(port);
	if err != nil {
		return ipAddr, errors.New(strings.Join([]string{"Error:", port, "is invalid port number"}, " "))
	}
	if (p <= 1024) {
		return ipAddr, errors.New(strings.Join([]string{"Error:", "Port", port, "must greater than", "1024"}, " "))
	}
	return ipAddr, nil 
}

func (config *ConfigType) Save() error {
	if _, err := CheckIPAddress(config.IpAddress, config.Port); err != nil {
		return err
	}
	data, err := json.Marshal(config);
	if err != nil {
		log.Panic(err)
	}
	f, err := os.OpenFile(filepath.FromSlash(fmt.Sprintf("%s%s%s", RootPath, "/", "config.json")), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Panic(err)
	}
	if _, err := f.Write(data); err != nil {
		log.Panic(err)
	}
	if err := f.Close(); err != nil {
		log.Error(err)
	}
	return nil
}

var (
	RootPath string
	Config ConfigType = ConfigType{IpAddress: "",Port: ""}
)

func init() {
	ex, err := os.Executable()
	if err != nil {
		log.Panic(err)
	}
	RootPath = filepath.Dir(ex)
	configPath := filepath.FromSlash(fmt.Sprintf("%s%s%s", RootPath, "/", "config.json"))
	log.Info("Config Path: " + configPath)
	if data, err := os.ReadFile(configPath); err != nil {
		log.Error(err)
		json := []byte("{}")
		if err := os.WriteFile(configPath, json, 0644); err != nil {
			log.Panic(err)
		} else {
			if _, err := os.ReadFile(configPath); err != nil {
				log.Panic(err)
			}
		}
	} else {
		if err := json.Unmarshal(data, &Config); err != nil {
			log.Panic(err)
		}
	}
}
