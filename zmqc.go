/*
 * Copyright (c) 2012, Joshua Foster
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above copyright
 *       notice, this list of conditions and the following disclaimer in the
 *       documentation and/or other materials provided with the distribution.
 *     * Neither the name of the <organization> nor the
 *       names of its contributors may be used to endorse or promote products
 *       derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package main

import (
	"bufio"
	"fmt"
	flag "github.com/droundy/goopt"
	zmq "github.com/pebbe/zmq4"
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

	var send, recv bool
	skip := false

	var socket *zmq.Socket
	switch *socket_type {
	case "PUSH":
		socket, _ = zmq.NewSocket(zmq.PUSH)
		send = true
	case "PULL":
		socket, _ = zmq.NewSocket(zmq.PULL)
		recv = true
	case "PUB":
		socket, _ = zmq.NewSocket(zmq.PUB)
		send = true
	case "SUB":
		socket, _ = zmq.NewSocket(zmq.SUB)
		recv = true
		if len(*subscriptions) == 0 {
			socket.SetSubscribe("")
		}
		for _, subscription := range *subscriptions {
			socket.SetSubscribe(subscription)
		}
	case "REQ":
		socket, _ = zmq.NewSocket(zmq.REQ)
		send = true
		recv = true
	case "REP":
		socket, _ = zmq.NewSocket(zmq.REP)
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
			socket.SendBytes([]byte(line), 0)
		}
		if recv {
			data, _ := socket.RecvBytes(0)
			writer.Write(data)
			writer.Flush()
		}
		if skip {
			skip = false
		}
	}

	fmt.Println("finished", *number)
}
