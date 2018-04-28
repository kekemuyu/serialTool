package controllers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"github.com/jacobsa/go-serial/serial"
)

var ReceiveType = 0
var Serials map[string]io.ReadWriteCloser

var ReceiverWs = make(map[string]*websocket.Conn, 100)

func init() {
	Serials = make(map[string]io.ReadWriteCloser)

	go func() {
		tmp := make([]byte, 1024)
		for {
			for k, v := range Serials {
				if n, err := v.Read(tmp); (n > 0) && (err == nil) {
					if v, ok := ReceiverWs[k]; ok == true {
						data := string(tmp[:n])
						if ReceiveType == 1 {
							data = hex.EncodeToString(tmp[:n])
						}
						msg := map[string]interface{}{"code": 0, "message": "serial receive success", "id": k, "data": data}
						err := v.WriteJSON(msg)
						if err != nil {
							delete(ReceiverWs, k)
							log.Printf("client.WriteJSON error: %v", err)
							v.Close()
						}
					}
				}
			}
		}
	}()
}

type SerialOpenController struct {
	beego.Controller
}

type SerialCloseController struct {
	beego.Controller
}

type SerialSendController struct {
	beego.Controller
}

type SerialReceiveController struct {
	beego.Controller
}

type SerialReceiveTypeController struct {
	beego.Controller
}

func (s *SerialOpenController) Post() {
	id := s.Ctx.Input.Param(":id")
	fmt.Println(id)
	if _, ok := Serials[id]; ok == true {
		s.Data["json"] = map[string]interface{}{"code": 1, "message": "id error"}
		s.ServeJSON()
		return
	}

	var opt serial.OpenOptions
	var err error
	if err = json.Unmarshal(s.Ctx.Input.RequestBody, &opt); err != nil {
		s.Data["json"] = map[string]interface{}{"code": 1, "message": "serial open fail"}
		s.ServeJSON()
		return
	}

	var tserial io.ReadWriteCloser
	tserial, err = serial.Open(opt)
	if err != nil {
		fmt.Println(err)
		s.Data["json"] = map[string]interface{}{"code": 1, "message": "serial open fail"}
		s.ServeJSON()
		return
	}
	Serials[id] = tserial
	s.Data["json"] = map[string]interface{}{"code": 0, "message": "serial open success"}
	s.ServeJSON()
}

func (s *SerialCloseController) Post() {
	id := s.Ctx.Input.Param(":id")
	if _, ok := Serials[id]; ok == false {
		s.Data["json"] = map[string]interface{}{"code": 1, "message": "id error"}
		s.ServeJSON()
		return
	}
	if err := Serials[id].Close(); err != nil {
		s.Data["json"] = map[string]interface{}{"code": 1, "message": "serial close fail"}
		s.ServeJSON()
		return
	}
	delete(Serials, id)
	s.Data["json"] = map[string]interface{}{"code": 1, "message": "serial close success"}
	s.ServeJSON()
}

func (s *SerialSendController) Post() {
	id := s.Ctx.Input.Param(":id")
	if _, ok := Serials[id]; ok == false {
		s.Data["json"] = map[string]interface{}{"code": 1, "message": "id error"}
		s.ServeJSON()
		return
	}
	var data = struct {
		Type int //发送数据类型
		Data string
	}{}
	if err := json.Unmarshal(s.Ctx.Input.RequestBody, &data); err != nil {
		s.Data["json"] = map[string]interface{}{"code": 1, "message": "serial send fail"}
		s.ServeJSON()
		return
	}

	var bdata []byte
	if data.Type == 1 {
		str := strings.Replace(data.Data, " ", "", -1)
		var err error
		bdata, err = hex.DecodeString(str)
		if err != nil {
			s.Data["json"] = map[string]interface{}{"code": 1, "message": "serial send fail,format error"}
			s.ServeJSON()
			return
		}
	} else {
		bdata = []byte(data.Data)
	}

	if _, err := Serials[id].Write(bdata); err != nil {
		delete(Serials, id)
		s.Data["json"] = map[string]interface{}{"code": 1, "message": "serial send fail"}
		s.ServeJSON()
		return
	}
	s.Data["json"] = map[string]interface{}{"code": 0, "message": "serial send success"}
	s.ServeJSON()
}

func (s *SerialReceiveController) Get() {
	fmt.Println("wesocket")
	id := s.GetString("id")
	ws, err := websocket.Upgrade(s.Ctx.ResponseWriter, s.Ctx.Request, nil, 1024, 1024)
	if err != nil {
		log.Fatal(err)
	}

	ReceiverWs[id] = ws
}

func (s *SerialReceiveTypeController) Get() {
	rtype := s.GetString("type")
	if rtype == "true" {
		ReceiveType = 1
	} else {
		ReceiveType = 0
	}

	s.Data["json"] = map[string]interface{}{"code": 0, "message": "type set ok"}
	s.ServeJSON()

}
