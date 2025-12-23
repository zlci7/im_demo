package main

import (
	"fmt"

	"net"
)

type User struct {
	Name string

	Addr string

	ChanUserMsg chan string

	conn net.Conn

	//关联的服务器

	server *Server
}

// 用户上线业务

func (u *User) Online() {

	u.server.mapLock.Lock()

	u.server.OnlineMap[u.Name] = u

	u.server.mapLock.Unlock()

	//广播用户上线

	u.server.BroadCast(u, "已上线")

}

// 用户下线业务

func (u *User) Offline() {

	u.server.mapLock.Lock()

	delete(u.server.OnlineMap, u.Name)

	u.server.mapLock.Unlock()

	//广播用户下线

	u.server.BroadCast(u, "已下线")

}

// 处理用户消息业务

func (u *User) DoMessage(msg string) {

	if msg == "who" {

		//查询当前在线用户

		u.server.mapLock.Lock()

		for _, user := range u.server.OnlineMap {

			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线...\n"

			u.ChanUserMsg <- onlineMsg

		}

		u.server.mapLock.Unlock()

	} else if len(msg) > 7 && msg[:7] == "rename|" {

		//修改用户名

		newName := msg[7:]

		//判断newName是否存在

		if _, isExist := u.server.OnlineMap[newName]; isExist {

			u.ChanUserMsg <- "当前用户名已经被使用\n"

		} else {

			u.server.mapLock.Lock()

			delete(u.server.OnlineMap, u.Name)

			u.server.OnlineMap[newName] = u

			u.server.mapLock.Unlock()

			u.Name = newName

			u.ChanUserMsg <- "您已经更新用户名:" + u.Name + "\n"

		}

	} else {

		//广播消息

		u.server.BroadCast(u, msg)

	}

}

// 创建一个用户的API

func NewUser(conn net.Conn, server *Server) *User {

	userAddr := conn.RemoteAddr().String()

	user := &User{

		Name: userAddr,

		Addr: userAddr,

		ChanUserMsg: make(chan string),

		conn: conn,

		server: server,
	}

	//启动监听线程

	go user.ListenUserMsg()

	return user

}

// 监听用户消息的channel的goroutine，一旦有消息就发送给客户端

func (u *User) ListenUserMsg() {

	defer fmt.Println("user.ListenUserMsg exit")

	for {

		msg, ok := <-u.ChanUserMsg

		if !ok {

			break

		}

		u.conn.Write([]byte(msg + "\n"))

	}

}
