package controller

import (
	"fmt"
	"server/internal/dto"
	"server/internal/service"
	"server/pkg/db"
	"server/pkg/msgpack"
	"server/pkg/utils"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

type PingController struct {
	znet.BaseRouter
	msgService *service.MessageService
}

func NewPingController() *PingController {
	return &PingController{
		msgService: service.NewMessageService(db.DB),
	}
}

func (p *PingController) Handle(request ziface.IRequest) {
	data := request.GetData()

	var msg dto.PingMessage

	if err := msgpack.Unmarshal(data, &msg); err != nil {
		fmt.Println(err)
	} else {
		jsonStr := utils.ToJson(msg)
		fmt.Println(jsonStr)
	}

	p.msgService.CreateMessage(string(data))

	msgID := request.GetMsgID()
	_ = request.GetConnection().SendMsg(msgID, data)
}

func (p *PingController) GetMessageService() *service.MessageService {
	return p.msgService
}
