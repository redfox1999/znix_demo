package main

import (
	"fmt"
	"runtime"
	"server/internal/router"
	"server/pkg/db"
	"server/pkg/logger"
	"sync"
	"time"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

var (
	connCount int32
	mutex     sync.Mutex
)

func OnConnStart(conn ziface.IConnection) {
	mutex.Lock()
	connCount++
	mutex.Unlock()
	//logger.Info("Client connectioned id %d addr: %s, current connections: %d", conn.GetConnID(), conn.RemoteAddrString(), connCount)
}

func OnConnStop(conn ziface.IConnection) {
	mutex.Lock()
	connCount--
	mutex.Unlock()
	//logger.Info("Client disconnected id %d addr: %s, current connections: %d", conn.GetConnID(), conn.RemoteAddrString(), connCount)
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

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			mutex.Lock()
			count := connCount
			mutex.Unlock()
			logger.Info("Clients: %d, goroutine count: %d", count, runtime.NumGoroutine())
		}
	}()

	router.InitRouter(s)
	logger.Info("Router initialized")
	s.Serve()
	logger.Info("Server Stoped")

}
