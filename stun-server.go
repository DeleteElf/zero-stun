package main

import (
	"flag"
	"github.com/deleteelf/goframework/utils/loghelper"
	"github.com/pion/stun/v3"
	"github.com/zeromicro/go-zero/core/conf"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"zero-stun/config"
)

var configFile = flag.String("f", "etc/server.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.ServerConfig
	conf.MustLoad(*configFile, &c)
	loghelper.GetLogManager().Init(loghelper.Debug)
	listenIp := "0.0.0.0"
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(listenIp), Port: c.Port})
	if err != nil {
		log.Fatalf("Failed to listen: %s", err)
	}

	log.Println("Listening on ==>" + listenIp + ":" + strconv.Itoa(c.Port))
	// 监听终止信号优雅关闭服务
	//defer listener.Close()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("收到关闭信号，正在停止 STUN 服务...")
		_ = listener.Close()
		os.Exit(0)
	}()
	for {
		buf := make([]byte, 1024)
		n, addr, err := listener.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Failed to read from UDP: %s", err)
			continue
		}
		msg := &stun.Message{Raw: buf[:n]}
		if err := msg.Decode(); err != nil {
			log.Printf("Failed to decode STUN message: %s", err)
			continue
		}
		username, success := msg.Attributes.Get(stun.AttrUsername)
		if success {
			log.Printf("Received STUN message success from %s =====> msg: %s ", addr, username)
		} else {
			log.Printf("Received STUN message failed from %s", addr)
			continue
		}
		// Handle the STUN message
		if msg.Type == stun.BindingRequest {
			// 构建响应，将客户端的公网/反射 IP 和端口写入 XORMappedAddress
			res, err := stun.Build(stun.NewTransactionIDSetter(msg.TransactionID), stun.BindingSuccess,
				&stun.XORMappedAddress{IP: addr.IP, Port: addr.Port},
			)
			if err != nil {
				log.Printf("[%s] 构建响应失败: %v", addr.String(), err)
				continue
			}
			// 将响应发回客户端
			if _, err := listener.WriteTo(res.Raw, addr); err != nil {
				log.Printf("[%s] 发送响应失败: %v", addr.String(), err)
			}
		}
	}
}
