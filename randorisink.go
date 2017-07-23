/* Task sink
*  Binds PULL socket to tcp://127.0.0.1:6666
*  Collects results from workers and fan via that socket
*/

package main

import (
	// "fmt"
	// "io"
	"log"
	"os"
	// "time"
	zmq "github.com/alecthomas/gozmq"
)

func main() {
	//create logfile with desired read/write permissions
	f, err := os.OpenFile("log-of-randori", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	// defer closing of context
	context, _ := zmq.NewContext()
	defer context.Close()

	//  Socket to receive messages on
	receiver, _ := context.NewSocket(zmq.PULL)
	defer receiver.Close()
	receiver.Bind("tcp://127.0.0.1:6666")

	// loop forever
	for {
		msgbytes, _ := receiver.Recv(0)
		log.Printf("%s\n", msgbytes)
	}

}
