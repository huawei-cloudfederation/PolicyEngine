package common

import (
	"sync"
	"net/http"
	"bytes"
	 "encoding/json"
	"io"
        "os"
	"log"
)

//Declare some structure that will eb common for both Anonymous and Gossiper modulesv
type DC struct {
	OutOfResource bool
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
}

type alldcs struct {
	Lck  sync.Mutex
	List map[string]*DC
}

//global consul config
type ConsulConfig struct {
	IsLeader    bool
	DCEndpoint  string
	StorePreFix string
	DCName      string
}

type PErequest struct{
        UnSupress bool
}

type SetThreshhold struct{
        Threshhold  int
}

//Declare somecommon types that will be used accorss the goroutines
var (
	ALLDCs             alldcs    //The data structure that stores all the Datacenter information
	ThisDCName         string    //This DataCenter's Name
	ResourceThresold   int       //Threshold value of any resource (CPU, MEM or Disk) after which we need to broadcast OOR
	TriggerPolicyCh    chan bool //Polcy Engine will listen in this Channel
	GossiperIp       string
)


func init() {
	TriggerPolicyCh = make(chan bool)
	ALLDCs.List = make(map[string]*DC)
	ResourceThresold = 100
}

func UnSupress(unsupress bool){
	var resp PErequest
	resp.UnSupress = unsupress

         data := new(bytes.Buffer)
	 err := json.NewEncoder(data).Encode(resp)
	 if err != nil {
		 log.Println("Error Marshalling the response")
		return
	}

        url := "http://" + GossiperIp + ":8080/v1/UNSUPRESS"
        res, _ := http.Post(url, "application/json; charset=utf-8",data)
        io.Copy(os.Stdout, res.Body)
}

func ThreshholdCh(threshhold int){
        var resp SetThreshhold 
        resp.Threshhold = threshhold

         data := new(bytes.Buffer)
         err := json.NewEncoder(data).Encode(resp)
	 if err != nil {
                 log.Println("Error Marshalling the response")
                return
        }

        url := "http://" + GossiperIp + ":8080/v1/THRESHHOLD/"
        res, _ := http.Post(url, "application/json; charset=utf-8",data)
        io.Copy(os.Stdout, res.Body)
}
