package policyengine

import (
	"log"
	"sync"
	 "encoding/json"
	 "io/ioutil"

	"fmt"
	"../common"
	"../consulib"
	"net/http"
	"time"
)

type dcData map[string]*common.DC

var (
	dcDataList dcData
)

type PE struct {
	*consulib.ConsulHandle
	policy           *Policy
	Current_DS_Index uint64

	Lck sync.Mutex
}


func NewPE(config *common.ConsulConfig) *PE {
	newPE := &PE{Current_DS_Index: 0}

	var ok bool
	newPE.ConsulHandle, ok = consulib.NewConsulHandle(config)
	if !ok {
		log.Println("NewPe: Error in creating a consulib Handle")
		return nil
	}

	return newPE

}

// ApplyNewPolicy
// Only
//
func (this *PE) ApplyNewPolicy() {

	for {
		//<-common.
		//Apply the this.Policy and take a decision
		fmt.Println("I am inside ApplyNewPolicy()\n")

		<-common.TriggerPolicyCh

		log.Println("ApplyNewPolicy: Called")

		if this.policy == nil {
			log.Println("ApplyNewPolicy: Cannot apply new policy since policy is nil")
			continue
		}
		this.Lck.Lock()
		ok := this.policy.TakeDecision()
		if ok != true {
			log.Println("ApplyNewPolicy: TakeDecision on new policy ", this.policy.Name, " failed")
		}
		this.Lck.Unlock()

	}
}

// populate policy
// read and update policy
// if leader repllicate
//
func (this *PE) UpdatePolicyFromDS(config *common.ConsulConfig) {
	for {

		log.Println("UpdatePolicyFromDS: called")

		data, resultingIndex, ok := this.WatchStore(this.Current_DS_Index)

		if ok && (this.Current_DS_Index < resultingIndex) {

			if config.IsLeader == true {
				err := this.ReplicateStore(data)
				if err != nil {
					log.Fatalln("UpdatePolicyFromDS: Data replication failed", err)
				}
			}
			this.Lck.Lock()
			//set the new ModifiedIndex
			this.Current_DS_Index = resultingIndex
			for _, value := range data.KVPairs { //only one policy will be passed since we store only one currently in our PE
				log.Println("UpdatePolicyFromDS: key and value ", string(value.Key), string(value.Value))
				newpolicy, err := this.ProcessNewPolicy(value.Key, value.Value)
				if err != nil {
					log.Println("UpdatePolicyFromDS: ProcessNewPolicy failes to ", err)
				}
				this.policy = newpolicy

			}
			this.Lck.Unlock()
		}
		//time.Sleep(5 * time.Second)
	}

}

// BootStrapPolicy is used to booststap policy
//
func (this *PE) BootStrapPolicy(config *common.ConsulConfig) {

	log.Println("BootStrapPolicy: Called")

	data, resultingIndex, ok := this.WatchStore(this.Current_DS_Index)
		fmt.Println("the dc is Leader ",config.IsLeader,this.Current_DS_Index,resultingIndex,ok)

	if ok && (this.Current_DS_Index < resultingIndex) {

		if config.IsLeader == true {
			err := this.ReplicateStore(data)
			if err != nil {
				log.Fatalln("BootStrapPolicy: Data replication failed", err)
			}
			//Since this gosspier is the leader he will unsupress the framewokrs
			log.Println("BootStrapPolicy: calling the Unsupress")
			common.UnSupress(true)
		}
		//set the new ModifiedIndex
		this.Current_DS_Index = resultingIndex
		for _, value := range data.KVPairs {
			log.Println("BootStrapPolicy: key and value ", string(value.Key), string(value.Value))
			newpolicy, err := this.ProcessNewPolicy(value.Key, value.Value)
			if err != nil {
				log.Println("BootStrapPolicy: ProcessNewPolicy failes to ", err)
			}
			this.policy = newpolicy
			//parse and keep the policy dont take decision now

		}

	} else {
		log.Fatalln("BootStrapPolicy: Failed to read from the store. Aborting")
	}

	log.Println("BootStrapPolicy: returning from BootStrapPolicy")

}
//Entry point for the policy engine
func Run(config *common.ConsulConfig) {

	log.Println("Run: PolicyEngine run called")

	pe := NewPE(config)
	go pe.GetAllDCdata() 

	pe.BootStrapPolicy(config)

	go pe.UpdatePolicyFromDS(config)
	go pe.ApplyNewPolicy()

}

func (this *PE) GetAllDCdata(){
        for {
        url := "http://"+ common.GossiperIp + ":8080/v1/ALLDCSTATUS"
        res, err := http.Get(url)

            if err != nil {
                panic(err.Error())
            }

            body, err := ioutil.ReadAll(res.Body)

            if err != nil {
                panic(err.Error())
            }

            var data []common.DC
            err=json.Unmarshal(body, &data)
		if err!=nil{
			fmt.Println("I am inGetAllDCdata ",err)
			return
		}

		fmt.Println("I am inGetAllDCdata ",data)
                common.ALLDCs.Lck.Lock()
                for i, v := range data{

			common.ALLDCs.List[v.Name]=&data[i]
		}
                fmt.Println("I am in main\n",common.ALLDCs.List)
                dcDataList = common.ALLDCs.List
                 common.ALLDCs.Lck.Unlock()
                time.Sleep(5 * time.Second)
        }

}
