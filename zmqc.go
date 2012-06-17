package main

import (
	"bufio"
	"fmt"
	zmq "github.com/alecthomas/gozmq"
	flag "github.com/droundy/goopt"
	"os"
)

var null = flag.Flag([]string{"-0"}, []string{}, "Separate messages on input/output should be ", "")

var mode = flag.Flag([]string{"-c", "--connect"},
	[]string{"-b", "--bind"},
	"Connect to the specified address(es).",
	"Bind to the specified address(es).")

var number = flag.Int([]string{"-n"}, -1, "Receive/send only NUM messages. By default, zmqc "+
	"lives forever in 'read' mode, or until the end of input "+
	"in 'write' mode.")

var socket_type = flag.Alternatives([]string{"-s"},
	[]string{"PUSH", "PULL", "PUB", "SUB", "REQ", "REP", "PAIR"},
	"Which type of socket to create. Must be one of 'PUSH', 'PULL', "+
		"'PUB', 'SUB', 'REQ', 'REP' or 'PAIR'. See `man zmq_socket` for an "+
		"explanation of the different types. 'DEALER' and 'ROUTER' sockets are "+
		"currently unsupported.")

var subscriptions = flag.Strings([]string{"--subscribe"}, "", "Subscribes to data matching")

func init() {
	flag.Version = "1.0"
	flag.Summary = "zmqc is a small but powerful command-line interface to " +
		"ZeroMQ. It allows you to create a socket of a given type, bind or " +
		"connect it to multiple addresses, set options on it, and receive or send " +
		"messages over it using standard I/O, in the shell or in scripts."
	flag.Author = "Joshua Foster"
}

func main() {
	flag.Parse(nil)

	address_list := flag.Args
	if len(address_list) == 0 {
		fmt.Println("No Addresses submitted")
		fmt.Println(flag.Help())
		return
	}

	context, _ := zmq.NewContext()

	var send, recv bool
	skip := false

	var socket zmq.Socket
	switch *socket_type {
	case "PUSH":
		socket, _ = context.NewSocket(zmq.PUSH)
		send = true
	case "PULL":
		socket, _ = context.NewSocket(zmq.PULL)
		recv = true
	case "PUB":
		socket, _ = context.NewSocket(zmq.PUB)
		send = true
	case "SUB":
		socket, _ = context.NewSocket(zmq.SUB)
		recv = true
		if len(*subscriptions) == 0 {
			socket.SetSockOptString(zmq.SUBSCRIBE, "")
		}
		for _, subscription := range *subscriptions {
			socket.SetSockOptString(zmq.SUBSCRIBE, subscription)
		}
	case "REQ":
		socket, _ = context.NewSocket(zmq.REQ)
		send = true
		recv = true
	case "REP":
		socket, _ = context.NewSocket(zmq.REP)
		send = true
		recv = true
		skip = true
	}
	defer socket.Close()

	// connect or bind
	if *mode {
		for _, address := range address_list {
			socket.Connect(address)
		}
	} else {
		for _, address := range address_list {
			socket.Bind(address)
		}
	}

	delim := byte('\n')
	if *null {
		fmt.Println("Setting delim to null")
		delim = byte(0x00)
	}

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	for i := 0; i < *number || *number == -1; i++ {
		if send && !skip {
			line, _ := reader.ReadBytes(delim)
			socket.Send([]byte(line), 0)
		}
		if recv {
			data, _ := socket.Recv(0)
			writer.Write(data)
			writer.Flush()
		}
		if skip {
			skip = false
		}
	}

	fmt.Println("finished", *number)
}
