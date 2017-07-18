package main

import (
	// "bufio"
	"bytes"
	"encoding/hex"
	"log"
	// "fmt"
	"net"
	// "os"
	"strings"
	"time"
	zmq "github.com/alecthomas/gozmq"
)

// telnet opcodes
// var WILL = []byte{251}
// var WONT = []byte{252}
// var DO   = []byte{253}
// var DONT = []byte{254}
var IAC  = []byte{255}
// 
func noyoudont(buf []byte, n int) (replace []byte) {
	// log.Printf("buf[0] = %v, this is real telnet", buf[0])
	// log.Printf(hex.EncodeToString(buf[:n])) // DEBUG
	// do -> dont ; will -> wont
	replace = bytes.Replace(buf[:n], []byte{255, 253}, []byte{255, 254}, -1)
	replace = bytes.Replace(replace, []byte{255, 251}, []byte{255, 252}, -1)
	// log.Printf(hex.EncodeToString(replace)) // DEBUG
	// log.Printf(hex.EncodeToString(buf)) // DEBUG
	replace = bytes.Replace(buf, []byte{255, 251}, []byte{255, 252}, -1)
	// log.Printf(hex.EncodeToString(replace)) // DEBUG

	return replace[:n]
}

func authcheck(ip string, username string, password string, out chan []byte) {
	var response []byte
	// log.Printf("%s:%s:%s", ip, username, password)
	conn, err := net.Dial("tcp", ip+":23")
	if err != nil {
		// log.Fatal("error conn failed: %s", err.Error())
		// out <- []byte(err.Error())
		out <- nil
		return
	}
	defer conn.Close()
	log.Printf("DEBUG: connecting to %s with %s:%s", ip, username, password)
	for {
		buf := make([]byte, 4096)
		// disconnect after *Millisecond if nothings in the buf
		err := conn.SetReadDeadline(time.Now().Add(4000 * time.Millisecond))
		if err != nil {
		// log.Println("SetReadDeadline failed:", err)
		break
		}
		n, _ := conn.Read(buf)
		if n > 0 {
			// Did I already mention TELNET is very ugly?
			if bytes.Equal([]byte{buf[0]}, IAC) {
				telnetnegotiation := noyoudont(buf, n)
				conn.Write(telnetnegotiation)
				// XXX interesting, but lots of 0's
				// to be investigated
				// response = append(response, buf...)

			} else {
				firstbuf := buf[:n]
				// log.Printf("firstbuf: %s", firstbuf)
				response = append(response, firstbuf...)

				// We cannot look for \n or EOF. Most if not all
				// login prompts end with a ':' but not
				// all.
				// So lets just trust what the malware gave us.
				// Wait. DID I JUST SAY THAT?!
				//
				// /bin/sh shite should probably
				// be filtered out upstream

				bannerresbuf, _ := conn.Read(buf)
				if bannerresbuf > 0 {
					secondbuf := buf[:bannerresbuf]
					// log.Printf("secondbuf: %s", secondbuf)
					response = append(response, secondbuf...)
				// write the first (username) field
					conn.Write([]byte(username+"\n"))

			}

				// commenting this out, as it creates uneven resposes
				usernamebuf, _ := conn.Read(buf)
				// Usually
				// it is just our username echo'd
				if usernamebuf > 0 {
					// thirdbuf := buf[:usernamebuf]
					// log.Printf("thirdbuf: %s", thirdbuf)

					//response = append(response, thirdbuf...)

			 }
				userresbuf, _ := conn.Read(buf)
				if userresbuf > 0 {
					fourthbuf := buf[:userresbuf]
					// log.Printf("fourthbuf: %s", fourthbuf)

					response = append(response, fourthbuf...)
					// write the password field
					conn.Write([]byte(password+"\n"))
			}
				passreadbuf, _ := conn.Read(buf)
				if passreadbuf > 0 {
					// Mostly this doesn't echo (password)
					// but just making sure we dont see it
					// (influences response length)
					//fifthbuf := buf[:passreadbuf]
					//log.Printf("fifthbuf: %s", fifthbuf)
					//response = append(response, fifthbuf...)

			}
					// conn.Write([]byte("\n"))
				passresponsebuf, _ := conn.Read(buf)
				if passresponsebuf > 0 {
					sixthbuf := buf[:passresponsebuf]
					// log.Printf("sixthbuf: %s", sixthbuf)

					response = append(response, sixthbuf...)
					// If this is real telnet, exit would be good
					// conn.Write([]byte("exit\n"))
			}
				loopreadbuf, _ := conn.Read(buf)
				if loopreadbuf > 0 {

					seventhbuf := buf[:passreadbuf]
					// log.Printf("seventhbuf: %s", seventhbuf)
					response = append(response, seventhbuf...)

			}
					// log.Printf("response: %s", hex.Dump(response)) // DEBUG
					// log.Printf("response: %s", hex.EncodeToString(response)) // DEBUG
					// response  = make([]byte, 0)
					out <- response
					// conn.Close()
					break
				}
			} else {
					// log.Printf("Hitting else: %s, %s", telnetline, response) // DEBUG
					out <- response
				        break
			}
		}
		// return
	}

func main() {
	out := make(chan []byte)

	context, _ := zmq.NewContext()
	socket, _ := context.NewSocket(zmq.SUB)

	defer context.Close()
	defer socket.Close()

	// var err error
	authsource := "login"
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
			log.Printf("telnet\t%s\t%s\t%s\t%s\t%d", ip, username, password, hex.EncodeToString(response),len(response))
		}
	}
}


