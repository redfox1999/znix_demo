package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

type Stats struct {
	success int64
	failed  int64
	mu      sync.Mutex
}

func (s *Stats) addSuccess(n int64) {
	s.mu.Lock()
	s.success += n
	s.mu.Unlock()
}

func (s *Stats) addFailed(n int64) {
	s.mu.Lock()
	s.failed += n
	s.mu.Unlock()
}

func sendMsg(conn net.Conn, msgID uint32, data []byte) error {
	header := make([]byte, 8)
	binary.BigEndian.PutUint32(header[0:4], msgID)
	binary.BigEndian.PutUint32(header[4:8], uint32(len(data)))
	_, err := conn.Write(append(header, data...))
	return err
}

func recvMsg(conn net.Conn) (uint32, []byte, error) {
	header := make([]byte, 8)
	if _, err := conn.Read(header); err != nil {
		return 0, nil, err
	}
	msgID := binary.BigEndian.Uint32(header[0:4])
	dataLen := binary.BigEndian.Uint32(header[4:8])
	body := make([]byte, dataLen)
	if _, err := conn.Read(body); err != nil {
		return 0, nil, err
	}
	return msgID, body, nil
}

func runClient(clientID int, host string, port int, repeat int, sem chan struct{}, stats *Stats, wg *sync.WaitGroup) {
	defer wg.Done()

	sem <- struct{}{}
	defer func() { <-sem }()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		stats.addFailed(int64(repeat))
		return
	}
	defer conn.Close()

	for i := 0; i < repeat; i++ {
		testData := map[string]interface{}{
			"arg1":   rand.Intn(99999) + 1,
			"arg2":   rand.Intn(99999) + 1,
			"result": 0,
		}

		data, err := msgpack.Marshal(testData)
		if err != nil {
			stats.addFailed(1)
			continue
		}

		if err := sendMsg(conn, 1001, data); err != nil {
			stats.addFailed(1)
			continue
		}

		if _, _, err := recvMsg(conn); err != nil {
			stats.addFailed(1)
		} else {
			stats.addSuccess(1)
		}
	}
}

func main() {
	serverHost := "127.0.0.1"
	serverPort := 8888
	totalConns := 100
	repeatPerConn := 1000

	fmt.Printf("🚀 启动压测: %d 并发连接, 每个连接请求 %d 次\n", totalConns, repeatPerConn)
	startTime := time.Now()

	stats := &Stats{}
	sem := make(chan struct{}, 1000)
	var wg sync.WaitGroup

	wg.Add(totalConns)
	for i := 0; i < totalConns; i++ {
		go runClient(i, serverHost, serverPort, repeatPerConn, sem, stats, &wg)
	}

	wg.Wait()
	duration := time.Since(startTime)

	qps := float64(stats.success) / duration.Seconds()

	fmt.Println("\n" + "=" + strings.Repeat("=", 39))
	fmt.Printf("🏁 压测报告\n")
	fmt.Printf("总耗时: %.2f 秒\n", duration.Seconds())
	fmt.Printf("成功次数: %d\n", stats.success)
	fmt.Printf("失败次数: %d\n", stats.failed)
	fmt.Printf("有效 QPS: %.2f\n", qps)
	fmt.Println("=" + strings.Repeat("=", 39))
}
