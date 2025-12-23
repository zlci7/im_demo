package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	//在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//消息广播的通道
	ChanServerMsg chan string
}

func NewServer(ip string, port int) *Server {

	server := &Server{
		Ip:            ip,
		Port:          port,
		OnlineMap:     make(map[string]*User),
		ChanServerMsg: make(chan string),
	}
	return server

}

func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := fmt.Sprintf("[%s]%s:%s\n", user.Addr, user.Name, msg)
	s.ChanServerMsg <- sendMsg

}

func (s *Server) ListenServerMsg() {
	for {
		msg := <-s.ChanServerMsg
		s.mapLock.Lock()
		//将消息发送给在线的用户
		for _, user := range s.OnlineMap {
			user.ChanUserMsg <- msg
		}
		s.mapLock.Unlock()
	}
}

func (s *Server) Handler(conn net.Conn) {
	//执行连结具体业务
	fmt.Println("连接建立成功")
	//当前用户上线
	user := NewUser(conn, s)
	//用户上线，将用户加入到在线用户列表中
	user.Online()
	//监听用户是否活跃的channel
	isLive := make(chan bool)

	//接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)

		for {
			n, err := conn.Read(buf)
			if n == 0 || err != nil {
				//用户下线
				user.Offline()
				return
			}
			if err != nil {
				fmt.Println("conn read error:", err)
				return
			}
			//提取用户消息（去除'\n'）
			msg := string(buf[:n-1])
			//广播消息
			user.DoMessage(msg)
			//用户的任意消息，代表当前用户是活跃的
			isLive <- true
		}
	}()

	for {
		select {
		case <-isLive:
			//当前用户是活跃的，应该重置定时器
		case <-time.After(time.Second * 600):
			//关闭连接
			conn.Close()
			//退出当前Handler
			return
		}

	}

}

func (s *Server) Start() {

	//1 socket listen
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("listen error:", err)
	}
	//4 close
	defer listen.Close()
	//启动监听msg的goroutine
	go s.ListenServerMsg()

	for {

		//2 accept
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen accept error:", err)
			continue
		}

		//3 handle
		go s.Handler(conn)
	}

}
