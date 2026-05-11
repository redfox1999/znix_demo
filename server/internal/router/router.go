package router

import (
	"server/internal/command"

	"github.com/aceld/zinx/ziface"

	"server/internal/controller"
)

func InitRouter(s ziface.IServer) {
	s.AddRouter(command.MsgIDPing, controller.NewPingController())
}
