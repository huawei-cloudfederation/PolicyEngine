package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"./common"
	"./policyengine"
)


type cConfig struct {
	ConsulConfig   common.ConsulConfig
}

func NewcConfig() cConfig {
	return cConfig{}
}


func ProcessConfFile(filename string, conf *cConfig) {

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
		config := NewcConfig()

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

		//Start the Policy Engine
		go policyengine.Run(&config.ConsulConfig)
		
		//wait for ever
		wait := make(chan struct{})
		<-wait

}
