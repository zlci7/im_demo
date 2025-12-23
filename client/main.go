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

	flag int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       99,
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

func (c *Client) menu() bool {
	var flag int
	fmt.Println("--------欢迎来到多人聊天系统--------")
	fmt.Println("1、公聊模式")
	fmt.Println("2、私聊模式")
	fmt.Println("3、更改用户名")
	fmt.Println("0、退出")

	fmt.Scanln(&flag)
	if flag < 0 || flag > 3 {
		return false
	}
	c.flag = flag
	return true
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "服务器ip地址(默认127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "服务器端口号(默认8888)")
}

func (c *Client) Run() {
	for c.flag != 0 {
		for c.menu() == false {
			fmt.Println(">>>>>>>>>请输入合法范围内的数字<<<<<<<<<")
		}
		switch c.flag {
		case 1:
			fmt.Println("公聊模式选择...")
		case 2:
			fmt.Println("私聊模式选择...")
		case 3:
			fmt.Println("更改用户名选择...")

		}
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

	client.Run()
}
