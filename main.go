package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"./common"
	"./policyengine"
//	"time"
//	"math/rand"
	"net/http"
)

type AllDC struct {
 	Name          string
        City          string
        Country       string
        Endpoint      string
	CPU           float64
        MEM           float64
        DISK          float64
        Ucpu          float64 //Remaining CPU
        Umem          float64 //Remaining Memory
        Udisk         float64 //Remaining Disk
        LastUpdate    int64   //Time stamp of current DC status
        LastOOR       int64   //Time stamp of when was the last OOR Happpend
        IsActiveDC    bool
	OutOfResource bool
}

type PolicyConfig struct {
	GossiperIP string
	ConsulConfig   common.ConsulConfig
}

func NewPolicyConfig() PolicyConfig {
	return PolicyConfig{}
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
		 
		fmt.Println(config)
		getAllDCdata()

		//Start the Policy Engine
		go policyengine.Run(&config.ConsulConfig)

		
		//wait for ever
		wait := make(chan struct{})
		<-wait

}

func getAllDCdata(){
	res, err := http.Get("http://54.201.4.103:8080/v1/ALLDCSTATUS")

	    if err != nil {
		panic(err.Error())
	    }

	    body, err := ioutil.ReadAll(res.Body)

	    if err != nil {
		panic(err.Error())
	    }

	    var data []AllDC
	    json.Unmarshal(body, &data)
		common.ALLDCs.Lck.Lock()
		var dc common.DC
		common.ALLDCs.List[dc.Name] = &dc
		for _, v := range data{

		dc.Name = v.Name
		dc.City = v.City
		dc.Country = v.Country
		dc.Endpoint = v.Endpoint
		dc.CPU = v.CPU
		dc.MEM = v.MEM
		dc.DISK = v.DISK
		dc.Ucpu = v.Ucpu
		dc.Umem = v.Umem
		dc.Udisk = v.Udisk
		dc.OutOfResource = v.OutOfResource
		dc.IsActiveDC = v.IsActiveDC
		dc.LastUpdate = v.LastUpdate
		dc.LastOOR = v.LastOOR
		}
		fmt.Println("I am in main\n",common.ALLDCs.List)

        	 common.ALLDCs.Lck.Unlock()

}
