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
	//logger.Info("Client connected", "conn_id", conn.GetConnID(), "addr", conn.RemoteAddrString(), "total", connCount)
}

func OnConnStop(conn ziface.IConnection) {
	mutex.Lock()
	connCount--
	mutex.Unlock()
	//logger.Info("Client disconnected", "conn_id", conn.GetConnID(), "addr", conn.RemoteAddrString(), "total", connCount)
}

func initDB() error {
	if err := db.LoadEnv(); err != nil {
		logger.Warn("Failed to load .env file", "error", err)
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
	logger.InitLoggerMulti("debug", []string{"stdout", "log/app.log"})

	if err := initDB(); err != nil {
		logger.Fatal("Database initialization failed", "error", err)
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
			logger.Info("Connection stats", "clients", count, "goroutines", runtime.NumGoroutine())
		}
	}()

	router.InitRouter(s)
	logger.Info("Router initialized")
	s.Serve()
	logger.Info("Server stopped")
}
