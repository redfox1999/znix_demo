package main

import (
	"fmt"
	"server/internal/router"
	"server/pkg/db"
	"server/pkg/logger"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

func OnConnStart(conn ziface.IConnection) {
	//logger.Info("Client connectioned id %d addr: %s", conn.GetConnID(), conn.RemoteAddrString())
	conn.SetProperty("userinfo", nil) // 可以放client 相关的数据
}

func OnConnStop(conn ziface.IConnection) {
	//logger.Info("Client disconnected id %d addr: %s", conn.GetConnID(), conn.RemoteAddrString())
}

func initDB() error {
	if err := db.LoadEnv(); err != nil {
		logger.WarnWithFields("Failed to load .env file", "error", err)
	}

	if err := db.InitDBFromEnv(); err != nil {
		if err := db.InitDB("sqlite3", "./data/app.db"); err != nil {
			return fmt.Errorf("failed to init database: %w", err)
		}
		logger.Info("Using default SQLite database")
	} else {
		logger.Info("Database initialized from environment")
	}

	return nil
}

func main() {
	// 初始化日志：同时输出到屏幕（彩色）和文件
	logger.InitLoggerMulti("debug", []string{"stdout", "log/app.log"})

	if err := initDB(); err != nil {
		logger.Fatal("Database initialization failed: %v", err)
		return
	}
	defer db.Close()

	s := znet.NewServer()
	s.SetOnConnStart(OnConnStart)
	s.SetOnConnStop(OnConnStop)

	router.InitRouter(s)
	logger.Info("Router initialized")
	s.Serve()
	logger.Info("Server Stoped")

}
