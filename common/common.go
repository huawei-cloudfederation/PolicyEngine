package common

import (
	"fmt"
	"sync"
	"net/http"
	"bytes"
	 "encoding/json"
	"io"
        "os"
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
	fmt.Printf("Initalizeing Common")

}

func UnSupress(data bool){
	var resp PErequest
	resp.UnSupress = data
	fmt.Println("resp.UnSupress is \n",resp.UnSupress)

         b := new(bytes.Buffer)
	 json.NewEncoder(b).Encode(resp)

        url := "http://" + GossiperIp + ":8080/v1/UNSUPRESS"
	fmt.Println("url is \n",url)
        res, _ := http.Post(url, "application/json; charset=utf-8",b)
        io.Copy(os.Stdout, res.Body)
}

func ThreshholdCh(data int){
        var resp SetThreshhold 
        resp.Threshhold = data
        fmt.Println("resp.Threshhold  is \n",resp.Threshhold )

         b := new(bytes.Buffer)
         json.NewEncoder(b).Encode(resp)

        url := "http://" + GossiperIp + ":8080/v1/THRESHHOLD/"
        fmt.Println("url is \n",url)
        res, _ := http.Post(url, "application/json; charset=utf-8",b)
        io.Copy(os.Stdout, res.Body)
}
