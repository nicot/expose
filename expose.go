package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

type options struct {
	server bool
	laddr  string
}

var usage = `Enter a port to expose (like 8080) or server to run expose in server mode.`
var server = `localhost:4567`

func flags(args []string) options {
	opts := options{}
	if len(args) != 2 {
		badArgs()
	}
	if args[1] == "serve" {
		opts.server = true
		return opts
	}
	/*
		we should verify arg1 is a valid listening address.
			n, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
			badArgs()
		}*/
	opts.laddr = args[1]
	return opts
}

func badArgs() {
	fmt.Println(usage)
	os.Exit(1)
}

func main() {
	opts := flags(os.Args)
	if opts.server {
		serve()
		return
	}
	connect(opts.laddr)
}

func log(s ...interface{}) {
	fmt.Println(s)
}

// serve listens on a port and proxies connections from the public
// internet to the program that `connects` back to it.
func serve() {
	log("serving")
	ln, err := net.Listen("tcp", server)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go newService(conn)
	}
}

// A new registered connection, that should have information proxied to it.
func newService(conn net.Conn) {
	ln, err := net.Listen("tcp", ":0") // Open up a new high port
	if err != nil {
		panic(err)
	}
}

// connect dials to a public facing server, then proxies packets from
// the upstream server to the local application at `laddr`.
func connect(laddr string) {
	log(laddr)
	serverConn, err := net.Dial("tcp", server)
	if err != nil {
		fmt.Println(err)
		return
	}
	serviceConn, err := net.Dial("tcp", laddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	go copy(serviceConn, serverConn)
	go copy(serverConn, serviceConn)
}

func copy(r io.Reader, w io.Writer) {
	buf := make([]byte, 2<<16)
	for {
		n, err := r.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			}
			panic(err)
		}
		buf = buf[:n]
		var read int
		for read < n {
			i, err := w.Write(buf)
			if err != nil {
				panic(err)
			}
			read += i
		}
	}
}
