package consulib

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"

	"../common"
)

// global map struct for all the dc list
// this will be populated only for the leader
type globalDCMap struct {
	*sync.Mutex
	//here string is the node name for now.
	// TODO: if two nodes on a WAN have the same name?
	//DCClientConnection map[string]*FederaConsulClient
	AvalibleDCInfo map[string]*api.AgentMember
}

type ConsulHandle struct {
	*api.Client
	*api.KV
	Name        string //this will have the dc name
	IsLeader    bool   //if this client is the leader
	DClist      *globalDCMap
	StorePrefix string
}

type KVData struct {
	api.KVPairs
	*api.QueryMeta
}

func NewglobalDCMap() *globalDCMap {
	return &globalDCMap{Mutex: new(sync.Mutex),
		AvalibleDCInfo: make(map[string]*api.AgentMember),
	}
}

func (this *ConsulHandle) PopulatetheGlobalDCMap() bool {

	log.Println("[INFO] PopulatetheGlobalDCMap: Called")

	clientAgent := this.Agent()
	//get wan members from agent

	localDCInfo, err := clientAgent.Self()

	if err != nil {
		log.Println("PopulatetheGlobalDCMap: Error retrieving the local Client info", err)
		return false
	}

	localDCName := localDCInfo["Config"]["Datacenter"].(string)
	log.Println("PopulatetheGlobalDCMap: Localagent name", localDCName)
	dcMembers, err := clientAgent.Members(true)
	if err != nil {
		log.Println("PopulatetheGlobalDCMap: Error list from members", err)
		return false
	}

	for index, member := range dcMembers {
		//Buf Fix
		newname := strings.Split(member.Name, ".")[1]
		if newname != localDCName {

			log.Println("PopulatetheGlobalDCMap: updating localDc dc name used", localDCName, member.Name, newname)
			this.DClist.AvalibleDCInfo[newname] = dcMembers[index]
		} else {
			log.Println("PopulatetheGlobalDCMap: Skipping localDc ", localDCName, member.Name)
		}
	}
	return true
}

//We need a global map so that it can be used by the policy

//this function to create the KV client object for the any DC's address passed.
func NewConsulHandle(config *common.ConsulConfig) (*ConsulHandle, bool) {

	var dcInfoMap *globalDCMap

	log.Println("[INFO]NewConsulHandle: called")

	defaultConfig := api.DefaultConfig()
	defaultConfig.Datacenter = config.DCName
	defaultConfig.Address = config.DCEndpoint

	log.Println("[INFO]NewConsulHandle:", defaultConfig)

	dcClient, err := api.NewClient(defaultConfig)
	if err != nil {
		log.Println("NewConsulHandle: Cannot connect to the Consul Server ", err)
		return nil, false
	}

	this := &ConsulHandle{Client: dcClient, Name: config.DCName, IsLeader: config.IsLeader, DClist: dcInfoMap, KV: dcClient.KV(), StorePrefix: config.StorePreFix}

	if config.IsLeader == true {
		//populate the DC maps
		this.DClist = NewglobalDCMap()
		ok := this.PopulatetheGlobalDCMap()
		if !ok {
			log.Println("NewConsulHandle: Unable to populate consul DC's list ")
			return nil, false
		}
	}

	return this, true
}

func (this *ConsulHandle) GetData(key string) []byte {
	data, Meta, err := this.KV.Get(key, nil)

	if err != nil {
		log.Println("GetData: Erorr in GetData ", err)
	}

	log.Println("GetData: Got data ", data, Meta)

	return data.Value
}

//WriteOptions can be used to override the datacenter

//key should have the prefix
func (this *ConsulHandle) PutData(key string, value []byte, dcName string) error {
	//writeMeta

	wop := &api.WriteOptions{Datacenter: dcName}
	_, err := this.Put(&api.KVPair{Key: key, Value: value}, wop)
	if err != nil {
		log.Println("PutData: Error Unable to put data to", this.Name, err, *this.Client)
		return err
	}
	return nil
}

func (this *ConsulHandle) GetList(waitIndex uint64) (*KVData, uint64, error) {

	log.Println("GetList waitindex ", waitIndex, this.StorePrefix)

	q := &api.QueryOptions{}
	q.WaitIndex = waitIndex
	q.WaitTime = time.Duration(time.Minute * 60)
	KeyValuelist, KeyValueMeta, err := this.List(this.StorePrefix+"/", q)
	if err != nil {
		log.Println("GetList: KV List failed", err, KeyValueMeta, this.Client)
		//we pass on the same index when it fails to fetch from the KV store
		// so that it can be re-tried
		return nil, waitIndex, err
	}
	newKVData := &KVData{}
	newKVData.KVPairs = KeyValuelist
	newKVData.QueryMeta = KeyValueMeta
	log.Println("GetList: The total len of data", len(newKVData.KVPairs))
	for _, val := range newKVData.KVPairs {

		log.Println("Data recevived from GetList", string(val.Key), string(val.Value), val)
	}

	log.Println("GetList: new waitindex ", KeyValueMeta.LastIndex)
	return newKVData, KeyValueMeta.LastIndex, nil
}
