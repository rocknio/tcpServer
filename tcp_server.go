package main

import (
	"flag"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"

	"go.uber.org/zap"
)

var (
	port = flag.Int("port", 9999, "Tcp Server Port")
	host = flag.String("bind", "0.0.0.0", "Bind IP Address")
)

// MainLogger 日志logger
var MainLogger *zap.Logger

// StatusChan 状态统计channel
var StatusChan chan string

// SendStatus 发送状态信息到统计携程
func SendStatus(statusType string) {
	StatusChan <- statusType
}

func startTicker() {
	ticker := time.NewTicker(1000 * time.Millisecond)
	go func() {
		for range ticker.C {
			PrintStatus()
		}
	}()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	InitLogger("tcpserver.log", "INFO")
	defer MainLogger.Sync()
	MainLogger.Debug("Logger init OK! FILE: tcpserver.log, LOGLEVEL: DEBUG")

	StatusChan = make(chan string, 1000000)
	go RecvStatus(StatusChan)

	var l net.Listener
	var err error
	l, err = net.Listen("tcp", *host+":"+strconv.Itoa(*port))
	if err != nil {
		MainLogger.Error("Error Listening: " + err.Error())
		os.Exit(1)
	}

	defer l.Close()
	MainLogger.Info("Listening on " + *host + ":" + strconv.Itoa(*port))

	isFirst := 0
	for {
		conn, err := l.Accept()
		if err != nil {
			MainLogger.Error("Error accepting: " + err.Error())
			continue
		}

		SendStatus("A")
		if isFirst != 1 {
			SendStatus("T")
			startTicker()
			isFirst = 1
		}
		go MsgHandler(conn)
	}
}
