package httplib

import (
	"encoding/json"
	"log"

	"github.com/astaxie/beego"

	"../common"
	"../policyengine"
)

type MainController struct {
	beego.Controller
}

type Triggerrequest struct{
        Policy bool
}

func (this *MainController) TriggerPolicy(){
        var data  Triggerrequest
        this.Data["Policy"] = this.Ctx.Input.Param(":Policy")

	log.Println("TriggerPolicyCh called :\n")
        err := json.Unmarshal(this.Ctx.Input.RequestBody,&data)
        if err != nil {
		log.Printf("Json Unmarshall error = %v", err)
        return
    }
	this.Ctx.Output.Body(this.Ctx.Input.RequestBody)

	if data.Policy == true {
		policyengine.GetAllDCdata()
		common.TriggerPolicyCh<-true
	}
}

func Run(config string) {

	log.Printf("Starting the HTTP server at port %s", config)

	beego.Run(":" + config)

}
