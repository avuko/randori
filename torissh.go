package main

import (
//	"crypto/rand"
//	"crypto/rsa"
//	"crypto/sha1"
//	"crypto/x509"
//	"encoding/binary"
//	"encoding/pem"
	"encoding/hex"
	// "net"
//	"github.com/kr/pty"
	"golang.org/x/crypto/ssh"
//	"io"
	"log"
	// "log/syslog"
//	"net"
//	"os"
//	"os/exec"
//	"strconv"
	"strings"
//	"sync"
//	"syscall"
//	"time"
//	"unsafe"
	zmq "github.com/alecthomas/gozmq"
)


// from https://raw.githubusercontent.com/golang-samples/cipher/master/crypto/rsa_keypair.go
// func buildkeys() (priv_pem []byte) {
// 	priv, err := rsa.GenerateKey(rand.Reader, 2014)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	err = priv.Validate()
// 	if err != nil {
// 		fmt.Println("Validation failed.", err)
// 	}
// 
// 	// Get der format. priv_der []byte
// 	priv_der := x509.MarshalPKCS1PrivateKey(priv)
// 
// 	// pem.Block
// 	// blk pem.Block
// 	priv_blk := pem.Block{
// 		Type:    "RSA PRIVATE KEY",
// 		Headers: nil,
// 		Bytes:   priv_der,
// 	}
// 
	// Resultant private key in PEM format.
	// priv_pem string
//	priv_pem = []byte(string(pem.EncodeToMemory(&priv_blk)))
//	log.Printf("info:\tGenerated a transient SSH private key")
//	return


func authcheck(ip string, username string, password string, out chan []byte) {
	// ssh client
	log.Printf("DEBUG: connecting to %s with %s:%s", ip, username, password)
	var response []byte
	for {
	sshConfig := &ssh.ClientConfig{
		User: string(username),
		Auth: []ssh.AuthMethod{ssh.Password(string(password))},
		ClientVersion: "SSH-2.0-OpenSSH_7.2p2 Ubuntu-4ubuntu2.2",
	}
	conn, err := ssh.Dial("tcp", ip+":22", sshConfig)
	if err != nil {
		out <- nil
		log.Printf("DEBUG: no service on %s", ip)
		return
	}

	defer conn.Close()
	response = conn.ServerVersion()
	log.Printf("DEBUG: success with %s, %s:%s", ip, username, password)
	out <- response
	conn.Close() // Kill connection after success.
	return
}
}



func main() {
	out := make(chan []byte)

	context, _ := zmq.NewContext()
	socket, _ := context.NewSocket(zmq.SUB)

	defer context.Close()
	defer socket.Close()

	// var err error
	authsource := "sshd"
	socket.SetSubscribe(authsource)
	socket.Connect("tcp://127.0.0.1:5556")
	log.Println("DEBUG: connected to zmq")
	for {
	datapt, _ := socket.Recv(0)
	rline := strings.Split(string(datapt), "\t")
	// probably skipping data source
	log.Printf("rline:%s, length:%d", rline[:], len(rline))
	ip, username, password := rline[1], rline[2], rline[3]
			go authcheck(ip[:], username[:], password[:], out)
			response := <-out
			// filter out zero length (error) responses
			if len(response) != 0 {
			log.Printf("sshd\t%s\t%s\t%s\t%s\t%d", ip, username, password, hex.EncodeToString(response),len(response))
		}
	}
}

