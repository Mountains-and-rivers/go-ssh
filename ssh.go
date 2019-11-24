package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"time"
)

type terminal struct {
	username string
	password string
	hostname string
	port     string
}

func main() {

	hostInfo := terminal{"moja", "Wgl,.2019", "47.111.77.29", "22"}
	commandList := []string{
		"sudo su",
		//"sudo -i",
		"whoami",
		"echo 'ssh output'",
		"exit", //退出切换后的用户
		"exit", //断开连接
	}

	//获取执行结果
	for {
		r := result(hostInfo, commandList)
		if *r == "" {
			fmt.Println("--------执行失败 重来！---------")
		} else {
			fmt.Println(*r)
			break
		}
	}
}

func result(hostInfo terminal, commands []string) *string {
	// SSH client config
	config := &ssh.ClientConfig{
		User: hostInfo.username,
		Auth: []ssh.AuthMethod{
			ssh.Password(hostInfo.password),
		},
		// Non-production only
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	// Connect to host
	client, err := ssh.Dial("tcp", hostInfo.hostname+":"+hostInfo.port, config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Create sesssion
	sess, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer sess.Close()

	// StdinPipe for commands
	stdin, err := sess.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	// 把结果保存在变量中
	var b bytes.Buffer
	sess.Stdout = &b
	sess.Stderr = &b

	// 终端输出结果
	//sess.Stdout = os.Stdout
	//sess.Stderr = os.Stderr

	// Start remote shell
	err = sess.Shell()
	if err != nil {
		log.Fatal(err)
	}

	// send the commands
	for _, cmd := range commands {
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Wait for sess to finish
	err = sess.Wait()
	if err != nil {
		log.Fatal(err)
	}
	// 打印结果
	//fmt.Println(b.String())
	ret := b.String()
	return &ret
}
