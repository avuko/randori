package main

import (
	"bytes"
	// "encoding/hex"
	"fmt"
	// "log"
	"net"
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

func noyoudont(buf []byte, n int) (replace []byte) {
	// log.Printf("buf[0] = %v, this is real telnet", buf[0]) // DEBUG
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

func authcheck(ip, username, password string) (response []byte) {
	var errmsg string
	conn, err := net.Dial("tcp", ip+":23")
	if err != nil {
		errmsg = fmt.Sprintf("ERROR:%s", err)
		response = []byte(errmsg)
		return response
	}
	defer conn.Close()
	for {
		buf := make([]byte, 4096)
		// disconnect after *Millisecond if nothing is in the buf
		err := conn.SetReadDeadline(time.Now().Add(4000 * time.Millisecond))
		if err != nil {
			errmsg = fmt.Sprintf("ERROR:%s", err)
			response = []byte(errmsg)
			return response
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

				// We cannot look for \n or EOF. Most
				// login prompts end with a ':'
				// ... but not all.
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
				return response
			}
		} else {
			errmsg = fmt.Sprintf("ERROR:n=0 for %s", ip)
			response = []byte(errmsg)
			return response
		}
	}

}

func main() {


	context, _ := zmq.NewContext()
	defer context.Close()

	//  Socket to receive messages on
	receiver, _ := context.NewSocket(zmq.PULL)
	defer receiver.Close()
	receiver.Connect("tcp://127.0.0.1:5023")

	//  Socket to send messages to task sink
	sender, _ := context.NewSocket(zmq.PUSH)
	defer sender.Close()
	sender.Connect("tcp://127.0.0.1:6666")

	for {
		msgbytes, _ := receiver.Recv(0)
		rline := strings.Split(string(msgbytes), "\t")
		ip, username, password := rline[1], rline[2], rline[3]
		result := authcheck(ip[:], username[:], password[:])
		// Remove newlines from telnet
		result = bytes.Replace(result, []byte("\r"), []byte(" "), -1)
		result = bytes.Replace(result, []byte("\n"), []byte(" "), -1)
	authcheckresult := fmt.Sprintf("TORITELNET: ip=%s username=%s password=%s result=%s", ip, username, password, result[:])
		sender.Send([]byte(authcheckresult), 0)
	}
}
