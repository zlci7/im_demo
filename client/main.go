package main

import (
	"flag"
	"fmt"
	"net"
)

var serverIp string
var serverPort int

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn

	flag string
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       "99",
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

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "服务器ip地址(默认127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "服务器端口号(默认8888)")
}

func (client *Client) DealResponse() {
	//一旦client.conn有数据，就直接copy到stdout标准输出上，永久阻塞监听
	for {
		buf := make([]byte, 1024)
		n, err := client.conn.Read(buf)
		if err != nil {
			fmt.Println("conn.Read err:", err)
			return
		}
		fmt.Print(string(buf[:n]))
	}
}

func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client.conn == nil {
		fmt.Println("连接服务器失败...")
		return
	}
	fmt.Println("连接服务器成功...")
	go client.DealResponse()

	client.Run()
}
