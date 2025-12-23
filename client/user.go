package main

import (
	"fmt"
	"time"
)

// 更新名称
func (c *Client) updateName() bool {
	fmt.Print(">>>请输入用户名:   ")
	fmt.Scanln(&c.Name)
	sendMsg := "rename|" + c.Name + "\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}

// 公聊模式
func (c *Client) PublicChat() {
	fmt.Println(">>>请输入聊天内容，exit退出.")
	var chatMsg string
	for {

		fmt.Scanln(&chatMsg)
		sendMsg := chatMsg + "\n"
		if chatMsg == "exit" {
			break

		}
		_, err := c.conn.Write([]byte(sendMsg))
		if err != nil {
			fmt.Println("conn.Write err:", err)
		}
	}

}

// 查询当前用户
func (c *Client) SelectUser() {
	sendMsg := "who\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}
}

// 私聊模式
func (c *Client) PrivateChat() {
	fmt.Print(">>>请输入聊天对象[用户名]: ")
	var remoteName string
	fmt.Scanln(&remoteName)
	var chatMsg string
	fmt.Print(">>>请输入聊天内容:")
	for {
		fmt.Scanln(&chatMsg)
		if chatMsg == "exit" {
			break
		}
		sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
		_, err := c.conn.Write([]byte(sendMsg))
		if err != nil {
			fmt.Println("conn.Write err:", err)
		}
	}

}

// 菜单
func (c *Client) menu() bool {
	var flag string
	fmt.Println("--------欢迎来到多人聊天系统--------")
	fmt.Println("1、公聊模式")
	fmt.Println("2、私聊模式")
	fmt.Println("3、更改用户名")
	fmt.Println("0、退出")

	fmt.Scanln(&flag)
	if flag < "0" || flag > "3" {
		return false
	}
	c.flag = flag
	return true
}

func (c *Client) Run() {
	for c.flag != "0" {
		for c.menu() == false {
			fmt.Println(">>>>>>>>>请输入合法范围内的数字<<<<<<<<<")
		}
		switch c.flag {
		case "1":
			fmt.Println("公聊模式选择...")
			c.PublicChat()
			time.Sleep(100 * time.Millisecond)
		case "2":
			fmt.Println("私聊模式选择...")
			c.SelectUser()
			c.PrivateChat()
			time.Sleep(100 * time.Millisecond)
		case "3":
			fmt.Println("更改用户名选择...")
			c.updateName()
			time.Sleep(100 * time.Millisecond)
		}
	}
}
