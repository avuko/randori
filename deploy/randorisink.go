/* Task sink
*  Binds PULL socket to tcp://127.0.0.1:6666
*  Collects results from workers and fan via that socket
*/

package main

import (
	// "fmt"
	// "io"
	"log"
        "log/syslog"
	// "os"
	// "time"
	zmq "github.com/alecthomas/gozmq"
)

func main() {
	//create logfile with desired read/write permissions
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "randorisink")
	if e == nil {
	log.SetFlags(0)
        log.SetOutput(logwriter)
    }

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
