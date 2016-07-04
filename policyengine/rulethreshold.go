package policyengine

import (
	"log"

	"../common"
	"strconv"
)

type RuleThreshold struct {
	ResourceLimit int `json:"ResourceLimit"`
}

func (this *RuleThreshold) ApplyRule(policydecision *PolicyDecision) bool {

	log.Println("ApplyRule: ApplyRule of the interface to RuleThreshold")
	//set the DC thershold
	//no chnage the dataset
	common.ResourceThresold = this.ResourceLimit
	dat := strconv.Itoa(common.ResourceThresold)
	common.ThreshholdCh(dat)
	log.Println("ApplyRule: RuleThreshold setting the Threshold", common.ResourceThresold)
	return true
}
