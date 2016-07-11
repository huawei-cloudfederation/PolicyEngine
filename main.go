package main

import (
	"./common"
	"./httplib"
	"./policyengine"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

type PolicyConfig struct {
	Name         string
	GossiperIP   string
	HTTPPort     string
	ConsulConfig common.ConsulConfig
}

func NewPolicyConfig() PolicyConfig {
	return PolicyConfig{
		HTTPPort: "8081",
	}
}

func ProcessConfFile(filename string, conf *PolicyConfig) {

	file_content, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatalf("Unable to read the config file %v", err)
	}

	err = json.Unmarshal(file_content, conf)

	if err != nil {
		log.Fatalf("unable to unmarshall the config file not a valid json err=%", err)
	}
}

func main() {

	log.Printf("The code just started")

	//Get the default Config populated just in case no config.json was supplied via comamnd line argument
	config := NewPolicyConfig()

	conffile := flag.String("config", "./config.json", "Supply the location of MrRedis configuration file")
	dummyConfig := flag.Bool("printDummyConfig", false, "IF you want to print the default(false) config")
	flag.Parse()

	if *dummyConfig == true {
		config_byte, err := json.MarshalIndent(config, " ", "  ")
		if err != nil {
			log.Printf("Error Marshalling the default config file %v", err)
			return
		}
		fmt.Printf("%s\n", string(config_byte))
		return

	}

	//Try to parse the config file
	ProcessConfFile(*conffile, &config)
	common.ThisDCName = config.Name
	common.GossiperIp = config.GossiperIP

	//start http server
	go httplib.Run(config.HTTPPort)

	//Start the Policy Engine
	go policyengine.Run(&config.ConsulConfig)

	//wait for ever
	wait := make(chan struct{})
	<-wait

}
