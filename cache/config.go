package cache

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// C 配置
type C struct {
	HostURL   string   `json:"hostURL"`
	Port      string   `json:"port"`
	WhiteList []string `json:"whiteList"`
}

// Config settings
var Config C

func init() {
	_, err := os.Open("config.json")
	if err != nil {
		f, _ := os.Create("config.json")
		c := C{}
		cb, _ := json.MarshalIndent(c, "", "  ")
		f.Write(cb)
		f.Close()
		log.Fatal("please complete the config.json")
	}
	data, _ := ioutil.ReadFile("config.json")
	err = json.Unmarshal(data, &Config)
	if err != nil {
		log.Fatal("config file is invalide")
	}
}
