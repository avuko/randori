// Task sink
// Binds PULL socket to tcp://127.0.0.1:5558
// Collects results from workers via that socket
//
package main

import (
	// "fmt"
	zmq "github.com/alecthomas/gozmq"
	// "time"
	"bytes"
	"log"
	"os"
	// "io"
)

func main() {
	//create your file with desired read/write permissions
	f, err := os.OpenFile("log-of-randori", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}   

	//defer to close when you're done with it, not because you think it's idiomatic!
	defer f.Close()

	//set output of logs to f
	log.SetOutput(f)

	context, _ := zmq.NewContext()
	defer context.Close()

	//  Socket to receive messages on
	receiver, _ := context.NewSocket(zmq.PULL)
	defer receiver.Close()
	receiver.Bind("tcp://127.0.0.1:6666")

	//  Wait for start of batch
	// msgbytes, _ := receiver.Recv(0)
	//    fmt.Println("Received Start Msg ", string(msgbytes))

	// loop forever
	for {
		msgbytes, _ := receiver.Recv(0)
		msgbytes = bytes.Replace(msgbytes, []byte("\r"), []byte(" "), -1)

		msgbytes = bytes.Replace(msgbytes, []byte("\n"), []byte(" "), -1)
		log.Printf("%s\n", msgbytes)
	}

}
