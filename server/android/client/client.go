package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	candy "github.com/dearcode/candy/server/android"
	"github.com/dearcode/candy/server/util/log"
)

func notice() {
	fmt.Println("---------------------------------")
	fmt.Println("1. 注册用户")
	fmt.Println("2. 登陆")
	fmt.Println("3. 注销")
	fmt.Println("4. 更新用户信息")
	fmt.Println("5. 获取用户信息")
	fmt.Println("6. 查找用户")
	fmt.Println("7. 添加好友")
	fmt.Println("0. 退出")
	fmt.Println("---------------------------------")
}

func register(c *candy.CandyClient, reader *bufio.Reader) {
	fmt.Println("================用户注册=======================")
	fmt.Println("请输入用户名:")
	data, _, _ := reader.ReadLine()
	userName := string(data)
	fmt.Println("请输入密码:")
	data, _, _ = reader.ReadLine()
	userPassword := string(data)

	id, err := c.Register(userName, userPassword)
	if err != nil {
		log.Errorf("Register error:%v", err)
		return
	}

	log.Debugf("Register success, userID:%v userName:%v userPassword:%v", id, userName, userPassword)

	fmt.Println("==============================================")
}

func login(c *candy.CandyClient, reader *bufio.Reader) {
	fmt.Println("================用户登陆=======================")
	fmt.Println("请输入用户名:")
	data, _, _ := reader.ReadLine()
	userName := string(data)
	fmt.Println("请输入密码:")
	data, _, _ = reader.ReadLine()
	userPassword := string(data)

	id, err := c.Login(userName, userPassword)
	if err != nil {
		log.Errorf("Login error:%v", err)
		return
	}

	log.Debugf("Login success, userID:%v userName:%v userPassword:%v", id, userName, userPassword)
	fmt.Println("==============================================")
}

func logout(c *candy.CandyClient, reader *bufio.Reader) {
	fmt.Println("================注销=======================")
	fmt.Println("请输入用户名:")
	data, _, _ := reader.ReadLine()
	userName := string(data)

	err := c.Logout(userName)
	if err != nil {
		log.Errorf("Logout error:%v", err)
		return
	}

	log.Debugf("Logout success, userName:%v", userName)
	fmt.Println("==============================================")
}

func updateUserInfo(c *candy.CandyClient, reader *bufio.Reader) {
	fmt.Println("================更新用户信息==================")
	fmt.Println("请输入用户名:")
	data, _, _ := reader.ReadLine()
	userName := string(data)
	fmt.Println("请输入用户昵称：")
	data, _, _ = reader.ReadLine()
	nickName := string(data)

	id, err := c.UpdateUserInfo(userName, nickName, nil)
	if err != nil {
		log.Errorf("updateUserInfo error:%v", err)
		return
	}

	log.Debugf("updateUserInfo success, userName:%v nickName:%v userID:%v", userName, nickName, id)
	fmt.Println("==============================================")
}

func getUserInfo(c *candy.CandyClient, reader *bufio.Reader) {
	fmt.Println("================获取用户信息==================")
	fmt.Println("请输入用户名:")
	data, _, _ := reader.ReadLine()
	userName := string(data)

	user, err := c.GetUserInfo(userName)
	if err != nil {
		log.Errorf("getUserInfo error:%v", err)
		return
	}

	log.Debugf("getUserInfo success, userName:%v", userName)
	log.Debugf("user detail, ID:%v Name:%v NickName:%v Avatar:%v", user.ID, user.Name, user.NickName, user.Avatar)
	fmt.Println("==============================================")
}

func findUser(c *candy.CandyClient, reader *bufio.Reader) {
	fmt.Println("================查找用户==================")
	fmt.Println("请输入用户名:")
	data, _, _ := reader.ReadLine()
	userName := string(data)

	users, err := c.FindUser(userName)
	if err != nil {
		log.Errorf("findUser error:%v", err)
		return
	}

	log.Debugf("findUser success, userName:%v*", userName)
	for index, user := range users {
		log.Debugf("user:%d detail, ID:%v Name:%v NickName:%v Avatar:%v", index, user.ID, user.Name, user.NickName, user.Avatar)
	}

	fmt.Println("==============================================")
}

func addFriend(c *candy.CandyClient, reader *bufio.Reader) {
	fmt.Println("================添加好友==================")
	fmt.Println("请输入用户ID:")
	data, _, _ := reader.ReadLine()
	userID := string(data)

	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		log.Errorf("Parse int error:%v", err)
		return
	}

	err = c.AddFriend(id, true)
	if err != nil {
		log.Errorf("addFriend error:%v", err)
		return
	}

	log.Debugf("addFriend success, userID:%v", userID)
	fmt.Println("==============================================")
}

func main() {
	c := candy.NewCandyClient("127.0.0.1:9000")
	if err := c.Start(); err != nil {
		log.Errorf("start client error:%s", err.Error())
		return
	}

	running := true
	reader := bufio.NewReader(os.Stdin)
	for running {
		notice()
		data, _, _ := reader.ReadLine()
		command := string(data)
		if command == "" {
			continue
		}

		log.Debugf("command:%v", command)
		if command == "0" {
			running = false
		} else if command == "1" {
			register(c, reader)
		} else if command == "2" {
			login(c, reader)
		} else if command == "3" {
			logout(c, reader)
		} else if command == "4" {
			updateUserInfo(c, reader)
		} else if command == "5" {
			getUserInfo(c, reader)
		} else if command == "6" {
			findUser(c, reader)
		} else if command == "7" {
			addFriend(c, reader)
		} else {
			log.Errorf("未知命令")
		}
	}

}