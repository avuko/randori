/* Task ventilator
*  Binds PUSH to multiple sockets
*  Sends batch of tasks to workers via these socket:
*  5021 = ftp
*  5022 = ssh
*  5023 = telnet
*  5080 = apache
*  6666 = overall sink
*/

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	zmq "github.com/alecthomas/gozmq"
)

func main() {
    fifofile := os.Stdin
    stdininfo, _ := fifofile.Stat()
    context, _ := zmq.NewContext()
    defer context.Close()

    // Socket to send telnet messages through [login]
    telnetsender, _ := context.NewSocket(zmq.PUSH)
    defer telnetsender.Close()
    telnetsender.Bind("tcp://127.0.0.1:5023")

    // Socket to send ssh messages through [sshd]
    sshsender, _ := context.NewSocket(zmq.PUSH)
    defer sshsender.Close()
    sshsender.Bind("tcp://127.0.0.1:5022")

    // Socket to send result messages through
    sink, _ := context.NewSocket(zmq.PUSH)
    defer sink.Close()
    sink.Connect("tcp://127.0.0.1:6666")

    // fmt.Println("Sending tasks to workers...") // DEBUG

	// Check whether stdininfo has a Mode and is a CharDevice
	if (stdininfo.Mode() & os.ModeCharDevice != 0) {
		log.Println("The command is intended to work with pipes.")
		log.Fatal("So, go use a pipe.")

	// Check whether stdininfo has a Mode and is a NamedPipe
	} else if (stdininfo.Mode() & os.ModeNamedPipe != 0) {
		scanner := bufio.NewScanner(fifofile)
		for scanner.Scan() {
			// log.Println(scanner.Text()) // DEBUG
			rline := strings.Split(scanner.Text(), "\t")
			// log.Printf("%d", len(rline)) // DEBUG
			if len(rline) == 5 {
			authsource, ip, username, password := rline[1], rline[2], rline[3], rline[4]
			// log.Printf("%s\t%s\t%s\t%s", authsource, ip, username, password) //DEBUG
			msg := fmt.Sprintf("%s\t%s\t%s\t%s", authsource, ip[:], username[:], password[:])

			// divide over different worker types
			switch {
			case authsource == "login":
				// log.Printf("DEBUG: protocol: %s", authsource) // DEBUG
			telnetsender.Send([]byte(msg), 0)
			case authsource == "sshd":
				// log.Printf("DEBUG: protocol: %s", authsource) // DEBUG
			sshsender.Send([]byte(msg), 0)
			default:
				// log.Printf("DEBUG: protocol: %s", authsource)  // DEBUG
			}
			} else {
				log.Printf("DEBUG: length of input line is %d", len(rline))
			}
				}
			if err := scanner.Err(); err != nil {
			log.Fatal("Error reading standard input: %s", err)
			}
	}

}
