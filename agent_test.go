package diste

import (
	"fmt"
	"math/rand"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	// "reflect"
	"sync"
	"testing"
	"time"
)

func setupServer(t *testing.T) net.Conn {
	portNum := rand.Intn(54000) + 10000
	URI := fmt.Sprintf("localhost:%d", portNum)

	agent := NewAgent("localhost")

	agentService := &AgentService{
		Agent: agent,
		Mutex: &sync.RWMutex{},
	}

	server := rpc.NewServer()
	server.Register(agentService)

	go func() {
		listen, listenErr := net.Listen(TCP_CONN, URI)
		if listenErr != nil {
			t.Fatal(listenErr)
		}
		defer listen.Close()

		for {
			conn, AcceptErr := listen.Accept()
			if AcceptErr != nil {
				t.Fatal(AcceptErr)
			}

			go server.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}()
	// Give some time for the service to start
	time.Sleep(100 * time.Millisecond)

	conn, dialError := net.Dial(TCP_CONN, URI)
	if dialError != nil {
		t.Fatal(dialError)
	}
	return conn
}

func TestAgentServiceFakeError(t *testing.T) {
	conn := setupServer(t)
	defer conn.Close()

	var reply *ServiceResponse
	var args ServiceRequest

	c := jsonrpc.NewClient(conn)
	err := c.Call("AgentService.FakeError", args, &reply)
	if err == nil {
		t.Fatal("Did not throw an error")
	}
}

func TestAgentServicePong(t *testing.T) {
	conn := setupServer(t)
	defer conn.Close()

	var reply *ServiceResponse
	var args ServiceRequest

	c := jsonrpc.NewClient(conn)
	err := c.Call("AgentService.Ping", args, &reply)
	if err != nil {
		t.Fatal(err)
	}
	if reply.Result != "Pong" {
		t.Fatal("assert fail", reply, "Pong")
	}
}

func TestAgentServiceState(t *testing.T) {
	conn := setupServer(t)
	defer conn.Close()

	var reply *ServiceResponse
	var args ServiceRequest

	c := jsonrpc.NewClient(conn)
	err := c.Call("AgentService.State", args, &reply)
	if err != nil {
		t.Fatal(err)
	}
	if reply.Result != string(OFF) {
		t.Fatal("assert fail", reply, OFF)
	}
}

func TestAgentServiceStart(t *testing.T) {
	conn := setupServer(t)
	defer conn.Close()

	var reply *ServiceResponse
	var args ServiceRequest

	c := jsonrpc.NewClient(conn)
	err := c.Call("AgentService.Start", args, &reply)
	if err != nil {
		t.Fatal(err)
	}
	if reply.Result != string(START) {
		t.Fatal("assert fail", reply, START)
	}
}
