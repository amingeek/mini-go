package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	port = "4001"
	path = "http://localhost:4001"
)

var serverSingleton *Server

func getServer() *Server {
	if serverSingleton == nil {
		serverSingleton = NewServer(port)
		go serverSingleton.Start()
		time.Sleep(1000 * time.Millisecond)
	}
	return serverSingleton
}

func TestServerCreation(t *testing.T) {
	s := getServer()
	assert.NotNil(t, s)
}
