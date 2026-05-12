package controller

import (
	"server/internal/dto"
	"server/internal/service"
	"server/pkg/db"
	"server/pkg/logger"
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
	msgID := request.GetMsgID()
	data := request.GetData()

	var msg dto.PingMessage

	if err := msgpack.Unmarshal(data, &msg); err != nil || !msg.Validate() {
		logger.WarnWithFields("Failed to unmarshal or validate ping message", "error", err)
		request.GetConnection().Stop()
		return
	} else {
		jsonStr := utils.ToJson(msg)
		logger.Print(jsonStr)
	}

	p.msgService.CreateMessage(string(data))

	_ = request.GetConnection().SendMsg(msgID, data)
}

func (p *PingController) GetMessageService() *service.MessageService {
	return p.msgService
}
