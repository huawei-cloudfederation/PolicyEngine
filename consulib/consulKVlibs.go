package consulib

import "log"

// This functon should have all the necessary information to replicate the KV store across
// the federated KV's
func (this *ConsulHandle) WatchStore(waitIndex uint64) (*KVData, uint64, bool) {

	log.Println("[INFO] WatchStore: called")
	//read from local store
	data, res, err := this.GetList(waitIndex)
	if err == nil && data != nil {
		log.Println("[INFO] WatchStore: Data received [o/p] new-index,old-index,data", res, waitIndex, data.KVPairs)
	} else {
		log.Println("[INFO] WatchStore: KV store dosent exist", err)
		return nil, waitIndex, false
	}

	return data, res, true
}

// THis will watch on the pre-defined store prefix and return the kv pairs
func (this *ConsulHandle) ReplicateStore(data *KVData) error {

	log.Println("ReplicateStore: called")
	if this.IsLeader != true {
		log.Println("ReplicateStore: Consul leader can only call this functions")
		return nil
	}
	//This function should iterate over all the DC's and push the message
	for dcName := range this.DClist.AvalibleDCInfo {
		//update all the DC's KV store in seperate goroutines
		//dc := this.DClist.DCClientConnection[dcName]
		for _, kvpair := range data.KVPairs {
			err := this.PutData(kvpair.Key, kvpair.Value, dcName)
			if err != nil {
				log.Println("ReplicateStore: Failed to send data to", dcName, this.DClist.AvalibleDCInfo[dcName])
				return err

			}
		}

	}

	return nil

}
