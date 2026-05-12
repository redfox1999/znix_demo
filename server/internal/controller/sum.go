package controller

import (
	"server/internal/command"
	"server/internal/dto"
	"server/pkg/logger"
	"server/pkg/msgpack"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

type SumController struct {
	znet.BaseRouter
}

func NewSumController() *SumController {
	return &SumController{}
}

func (s *SumController) Handle(req ziface.IRequest) {
	data := req.GetMessage().GetData()

	var sumMsg dto.SumMessage

	if err := msgpack.Unmarshal(data, &sumMsg); err != nil || !sumMsg.Validate() {
		logger.Error("Failed to unmarshal or validate sum message", "error", err, "sumMsg", sumMsg)
		req.GetConnection().Stop()
		return
	}

	var sumRes dto.SumResponse
	sumRes.Result = *sumMsg.Arg1 + *sumMsg.Arg2
	sumRes.Arg1 = *sumMsg.Arg1
	sumRes.Arg2 = *sumMsg.Arg2

	data, err := msgpack.Marshal(sumRes)
	if err != nil {
		logger.Error("Marshal sum response failed", "error", err)
		return
	}
	_ = req.GetConnection().SendMsg(command.MsgIDSum, data)
}