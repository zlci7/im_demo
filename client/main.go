package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}

	//连接服务器
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error", err)
		return nil
	}

	//返回对象
	client.conn = conn
	return client
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "服务器ip地址(默认127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "服务器端口号(默认8888)")
}

func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client.conn == nil {
		fmt.Println("连接服务器失败...")
		return
	}
	fmt.Println("连接服务器成功...")

	select {}
}
