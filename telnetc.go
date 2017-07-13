package main

import (
	"bytes"
	"log"
	"net"
	"os"
	"bufio"
	"strings"
	"time"
	"encoding/hex"
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

func telnetcheck(telnetline string, out chan []byte) {
	var response []byte
	telnetinfo := strings.Split(telnetline, "\u2002")
	log.Printf("%T:%s", telnetinfo, telnetinfo)
	ip, username, password := telnetinfo[0], telnetinfo[1], telnetinfo[2]
	conn, err := net.Dial("tcp", ip+":23")
	if err != nil {
		// log.Fatal("error conn failed: %s", err.Error())
		out <- []byte(err.Error())
		return
	}
	defer conn.Close()
	for {
		buf := make([]byte, 4096)
		// disconnect after *Millisecond if nothings in the buf
		err := conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		if err != nil {
		log.Println("SetReadDeadline failed:", err)
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
					response = append(response, secondbuf...)
				// write the first (username) field
					conn.Write([]byte(username+"\n"))

			}

				usernamebuf, _ := conn.Read(buf)
				// Add to response, although usually
				// it is just our username echo'd
				if usernamebuf > 0 {

					response = append(response, buf[:usernamebuf]...)

			}
				userresbuf, _ := conn.Read(buf)
				if userresbuf > 0 {

					response = append(response, buf[:userresbuf]...)
					// write the password field
					conn.Write([]byte(password+"\n"))
			}
				passreadbuf, _ := conn.Read(buf)
				if passreadbuf > 0 {

					response = append(response, buf[:passreadbuf]...)

			}
				passreadbuf2, _ := conn.Read(buf)
				if passreadbuf2 > 0 {

					response = append(response, buf[:passreadbuf2]...)
					// If this is real telnet, exit would be good
					// conn.Write([]byte("exit\n"))
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
	//fifofile, err := os.OpenFile("login", os.O_RDONLY, 0600)
	fifofile := os.Stdin
	//if err != nil {
	//		panic(err)
	//	    }
	stdininfo, _ := fifofile.Stat()

	out := make(chan []byte)
	// Check whether stdininfo has a Mode and is a CharDevice
	if (stdininfo.Mode() & os.ModeCharDevice != 0) {
		log.Println("The command is intended to work with pipes.")
		log.Fatal("So, go use a pipe.")
	// Check whether stdininfo has a Mode and is a NamedPipe
	} else if (stdininfo.Mode() & os.ModeNamedPipe != 0) {
		scanner := bufio.NewScanner(fifofile)
		for scanner.Scan() {
			// log.Println(scanner.Text())
			go telnetcheck(scanner.Text(), out)
			log.Printf("telnet\u2002%s\u2002%s", scanner.Text(), hex.EncodeToString(<-out))
			// log.Printf("telnet\u2002%s\u2002%s", scanner.Text(), <-out)
				}
			if err := scanner.Err(); err != nil {
			log.Fatal("Error reading standard input: %s", err)
			}
	}
	}


