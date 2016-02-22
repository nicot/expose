package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/hashicorp/yamux"
	"io/ioutil"
)

func main() {
}

type host struct {
	port     string
	hostName string
}

func (h *host) String() string {
	return h.hostName + h.port
}

var defaultUpstream = host{
	port:     ":3800",
	hostName: "localhost",
}

func parseArgs(args []string) error {
	if len(os.Args) < 2 {
		return errNotEnoughArgs
	}
	if os.Args[1] == "serve" {
		return serve(defaultUpstream.port)
	}
	_, err := strconv.ParseUint(os.Args[1], 10, 16)
	if err != nil {
		return err
	}
	h := host{port: ":" + os.Args[1]}
	return expose(defaultUpstream, h)
}

func serve(port string) error {
	l, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go handle(upstreamServe(conn))
	}

	return nil
}

func expose(upstream host, downstream host) error {
	conn, err := net.Dial("tcp", upstream.String())
	if err != nil {
		return err
	}
	return downstreamServe(conn, downstream)
}

func downstreamServe(conn net.Conn, downstream host) error {
	session, err := yamux.Server(conn, nil)
	if err != nil {
		return err
	}

	control, err := session.Accept()
	if err != nil {
		return err
	}

	laddr, err := ioutil.ReadAll(control)
	fmt.Println(string(laddr))
	if err != nil {
		return err
	}

	for {
		incoming, err := session.Accept()
		if err != nil {
			return err
		}

		outgoing, err := net.Dial("tcp", downstream.String())
		if err != nil {
			return err
		}

		go handle(proxy(incoming, outgoing))
	}

	return nil
}

func handle(err error) {
	fmt.Fprintln(os.Stderr, err)
}

func upstreamServe(conn net.Conn) error {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return err
	}
	session, err := yamux.Client(conn, nil)
	if err != nil {
		return err
	}

	control, err := session.Open()
	if err != nil {
		return err
	}

	laddr := l.Addr().String()
	if _, err := control.Write([]byte(laddr)); err != nil {
		return err
	}

	if err := control.Close(); err != nil {
		return err
	}

	for {
		incoming, err := l.Accept()
		if err != nil {
			return err
		}

		downstream, err := session.Open()
		if err != nil {
			return err
		}
		go handle(proxy(incoming, downstream))
	}

	return nil
}

func proxy(incoming, downstream net.Conn) error {
	// copy from incoming to downstream and vice versa
	return nil
}
