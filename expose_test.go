package main

import (
	"bytes"
	"net"
	"testing"
)

type proxy struct {
	// A server will probably be provided by us, and should have a
	// publically accessible IP.
	server net.Listener
	up     net.Conn

	// A service is a networked application that is presumably behind
	// a firewall or nat. An example would be a developer running a
	// rails server.
	service net.Listener
	down    net.Conn
}

func setup(t *testing.T) proxy {
	var p proxy
	var err error
	serverPort := ":4567"
	p.server, err = net.Listen("tcp", serverPort)
	if err != nil {
		t.Fatal(err)
	}

	servicePort := ":8080"
	p.service, err = net.Listen("tcp", servicePort)
	if err != nil {
		t.Fatal(err)
	}

	go connect(servicePort)

	p.up, err = p.server.Accept()
	if err != nil {
		t.Fatal(err)
	}
	p.down, err = p.service.Accept()
	if err != nil {
		t.Fatal(err)
	}

	return p
}

func (p proxy) Close() {
	p.up.Close()
	p.server.Close()

	p.down.Close()
	p.service.Close()
}

func TestUpstream(t *testing.T) {
	p := setup(t)
	defer p.Close()

	data := []byte("A little test string")
	p.up.Write(data)

	b := make([]byte, len(data))
	_, err := p.down.Read(b)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, b) {
		t.Errorf("Expected %x, got %x", data, b)
	}
}

func TestDownstream(t *testing.T) {
	p := setup(t)
	defer p.Close()

	data := []byte("hello world")
	p.down.Write(data)

	b := make([]byte, len(data))
	_, err := p.up.Read(b)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, b) {
		t.Errorf("Expected %x, got %x", data, b)
	}
}

func TestServe(t *testing.T) {
	servicePort := ":8080"
	ln, err := net.Listen("tcp", servicePort)
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	go serve()
	go connect(servicePort)

	conn, err := ln.Accept()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client, err := net.Dial("tcp", server)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	data := []byte("GET A new era")
	_, err = client.Write(data)
	if err != nil {
		t.Fatal(err)
	}

	b := make([]byte, len(data))
	_, err = conn.Read(b)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, b) {
		t.Errorf("Expected %x, got %x", data, b)
	}
}
