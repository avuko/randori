package main

import (
	// "bytes"
	"fmt"
	"log"
	// "net"
	"os"
	"bufio"
	"strings"
	"time"
	// "encoding/hex"
	zmq "github.com/alecthomas/gozmq"
)

func main() {
	//fifofile, err := os.OpenFile("login", os.O_RDONLY, 0600)
	fifofile := os.Stdin
	//if err != nil {
	//		panic(err)
	//	    }
	stdininfo, _ := fifofile.Stat()
	// out := make(chan []byte)

	// ZMQ startup
	context, _ := zmq.NewContext()
	socket, _ := context.NewSocket(zmq.PUB)
	defer context.Close()
	defer socket.Close()
	socket.Bind("tcp://127.0.0.1:5556")
	time.Sleep(time.Second * 5);
	// Check whether stdininfo has a Mode and is a CharDevice
	if (stdininfo.Mode() & os.ModeCharDevice != 0) {
		log.Println("The command is intended to work with pipes.")
		log.Fatal("So, go use a pipe.")
	// Check whether stdininfo has a Mode and is a NamedPipe
	} else if (stdininfo.Mode() & os.ModeNamedPipe != 0) {
		scanner := bufio.NewScanner(fifofile)
		for scanner.Scan() {
			// log.Println(scanner.Text())
			rline := strings.Split(scanner.Text(), "\t")
			log.Printf("%d", len(rline))
			if len(rline) == 5 {
			authsource, ip, username, password := rline[1], rline[2], rline[3], rline[4]
			log.Printf("%s\t%s\t%s\t%s", authsource, ip, username, password)
			msg := fmt.Sprintf("%s\t%s\t%s\t%s", authsource, ip[:], username[:], password[:])
			socket.Send([]byte(msg), 0)
			}
				}
			if err := scanner.Err(); err != nil {
			log.Fatal("Error reading standard input: %s", err)
			}
	}
	time.Sleep(time.Second * 5);
	}


