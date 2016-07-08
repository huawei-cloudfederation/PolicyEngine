package httplib

import (
	"encoding/json"
	"log"

	"github.com/astaxie/beego"

	"../common"
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

	log.Println("TriggerPolicyCh called in gossiper:\n")
        err := json.Unmarshal(this.Ctx.Input.RequestBody,&data)
           log.Println(this.Ctx.Input.RequestBody)
        if err != nil {
                this.Ctx.Output.Body(this.Ctx.Input.RequestBody)
                log.Println(this.Ctx.Input.RequestBody)
        return
    }
	this.Ctx.Output.Body(this.Ctx.Input.RequestBody)
                log.Println(string(this.Ctx.Input.RequestBody),"::",data)
	if data.Policy {
		common.TriggerPolicyCh<-true
	}
}

func Run(config string) {

	log.Printf("Starting the HTTP server at port %s", config)

	beego.Run(":" + config)

}
