package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	// "log"
	"strings"
	zmq "github.com/alecthomas/gozmq"
)

func authcheck(ip, username, password string) (response []byte) {
	var errmsg string
	sshConfig := &ssh.ClientConfig{
		User: string(username),
		Auth: []ssh.AuthMethod{ssh.Password(string(password))},
		ClientVersion: "SSH-2.0-OpenSSH_7.2p2 Ubuntu-4ubuntu2.2",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	connection, err := ssh.Dial("tcp", ip+":22", sshConfig)
	if err != nil {
		errmsg = fmt.Sprintf("ERROR:%s", err)
		response = []byte(errmsg)
		return response
	}
	response = connection.ServerVersion()
	connection.Close() // Kill connection after success.
	return response
}

func main() {

    context, _ := zmq.NewContext()
    defer context.Close()

    //  Socket to receive messages on
    receiver, _ := context.NewSocket(zmq.PULL)
    defer receiver.Close()
    receiver.Connect("tcp://127.0.0.1:5022")

    //  Socket to send messages to task sink
    sender, _ := context.NewSocket(zmq.PUSH)
    defer sender.Close()
    sender.Connect("tcp://127.0.0.1:6666")

	for {
	msgbytes, _ := receiver.Recv(0)
	rline := strings.Split(string(msgbytes), "\t")
	ip, username, password := rline[1], rline[2], rline[3]
	result := authcheck(ip[:], username[:], password[:])
	authcheckresult := fmt.Sprintf("TORISSH: ip=%s username=%s password=%s result=%s", ip, username, password, result[:])
	sender.Send([]byte(authcheckresult), 0)
	}

}
