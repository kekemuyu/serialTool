package routers

import (
	"serialTool/controllers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	beego.Router("/", &controllers.MainController{})
	beego.Router("/serial/open/:id", &controllers.SerialOpenController{})
	beego.Router("/serial/close/:id", &controllers.SerialCloseController{})
	beego.Router("/serial/send/:id", &controllers.SerialSendController{})
	beego.Router("/serial/receive", &controllers.SerialReceiveController{})
	beego.Router("/serial/receive/type", &controllers.SerialReceiveTypeController{})
}
