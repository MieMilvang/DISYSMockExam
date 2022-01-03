package main

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"

	h "github.com/MieMilvang/DISYSMockExam/HelperMethod"
	Proto "github.com/MieMilvang/DISYSMockExam/Proto"

	"google.golang.org/grpc"
)

const (
	SERVER_PORT     = 5000
	SERVER_LOG_FILE = "serverLog"
)

type Server struct {
	Proto.UnimplementedProtoServiceServer
	port        int
	latestValue h.Value
	arbiter     sync.Mutex
}
var (
	incrementedValue int = -1
)

func main() {
	//init
	initValue := h.Value{Value: -1, UserId: -1}
	serverPort := FindFreePort()
	if serverPort == -1 { // if no free port -1
		fmt.Printf("Can't start more than %v", h.MAX_REPLICAS)
		return
	}
	server := Server{port: serverPort, latestValue: initValue, arbiter: sync.Mutex{}}
	fmt.Printf("Succesfully got port: %v\n", server.port) // sanity checks

	listen(&server)
	fmt.Println("main has ended")
}

// get value grpc method logic
func (s *Server) GetValue(ctx context.Context, request *Proto.GetRequest) (*Proto.Value, error) {
	value := Proto.Value{CurrentValue: s.latestValue.Value, UserId: s.latestValue.UserId}
	return &value, nil
}

// set value grpc method logic


func (s *Server) Increment(ctx context.Context, request *Proto.SetRequest)(*Proto.Value, error){
	s.arbiter.Lock()
	//temp := s.latestValue
	incrementedValue++;
	//s.latestValue = h.Value{Value: request.GetRequestedValue(), UserId: request.GetUserId()}
	//msg := fmt.Sprintf("Updated the value: %v by %v to %v by user %v ", temp.Value, temp.UserId, s.latestValue.Value, s.latestValue.UserId)
	//h.Logger(msg, SERVER_LOG_FILE+strconv.Itoa(s.port))
	s.arbiter.Unlock()
	return &Proto.Value{CurrentValue: int64(incrementedValue), UserId: s.latestValue.UserId}, nil 
}

func (s *Server) JoinService(ctx context.Context, request *Proto.JoinRequest) (*Proto.Response, error) {
	userId := request.GetUserId()
	if userId == -1 {
		return &Proto.Response{Msg: "alive"}, nil
	} else {
		msg := fmt.Sprintf("Welcome to our marvelous service user: %v ", userId)
		return &Proto.Response{Msg: msg}, nil
	}
}

// connect to ports until a free port is found
func FindFreePort() int {
	for i := 1; i < (h.MAX_REPLICAS + 1); i++ {
		serverPort := SERVER_PORT + i
		_, status := h.ConnectToPort(serverPort)
		if status == "alive" {
			continue
		} else {
			return serverPort
		}
	}
	return -1
}

// start server service
func listen(s *Server) {

	//listen on port
	lis, err := net.Listen("tcp", "localhost:"+strconv.Itoa(s.port))
	h.CheckError(err, "server setup net.listen")
	defer lis.Close()

	// register server this is a blocking call
	grpcServer := grpc.NewServer()
	Proto.RegisterProtoServiceServer(grpcServer, s)
	errorMsg := grpcServer.Serve(lis)
	h.CheckError(errorMsg, "server listen register server service")
}
