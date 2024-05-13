package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	pb "poc-serialization/proto"
	"strings"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	pb.SerializedServiceServer
	parentData  parent
	state       state
	globalMutex sync.Mutex
	done        chan bool
	smaphore    chan bool
}

type state struct {
	counter int
}

type parent struct {
	childs []*child
}

type child struct {
	status  bool
	client  pb.SerializedServiceClient
	counter int
}

var (
	ports    = flag.String("ports", "50002,50003", "")
	mainPort = flag.Int("mainport", 50001, "")
)

func (s *server) ChildCC(ctx context.Context, in *pb.CC) (*pb.UnaryResponse, error) {
	return &pb.UnaryResponse{
		Ok: in.Counter == int32(s.state.counter),
	}, nil
}

func (s *server) ParentCC(ctx context.Context, in *pb.CC) (*pb.UnaryResponse, error) {
	s.globalMutex.Lock()
	defer s.globalMutex.Unlock()

	okCounter := new(int32)

	wg := sync.WaitGroup{}

	for _, chi := range s.parentData.childs {
		wg.Add(1)

		go func(chi *child, wg *sync.WaitGroup) {
			defer wg.Done()
			resp, _ := chi.client.ChildCC(ctx, &pb.CC{
				Counter: int32(chi.counter),
			})

			if resp.Ok {
				atomic.AddInt32(okCounter, 1)
			}
		}(chi, &wg)
	}

	wg.Wait()
	len := len(s.parentData.childs)
	ok := *okCounter == int32(len)

	if !ok {
		return nil, fmt.Errorf("")
	}

	return &pb.UnaryResponse{
		Ok: ok,
	}, nil
}

func (s *server) ParentUnaryCall(context context.Context, request *pb.UnaryRequest) (*pb.UnaryResponse, error) {
	return s.ParentUnaryCallV1(context, request)
}

func (s *server) ParentUnaryCallV1(context context.Context, request *pb.UnaryRequest) (*pb.UnaryResponse, error) {
	s.globalMutex.Lock()
	defer s.globalMutex.Unlock()

	leng := int32(len(s.parentData.childs))
	okCounter := new(int32)

	wg := sync.WaitGroup{}

	for _, chi := range s.parentData.childs {
		wg.Add(1)
		go func(chi *child) {
			defer wg.Done()
			resp, err := chi.client.UnaryCall(context, request)

			if resp != nil && resp.Ok && err == nil {
				atomic.AddInt32(okCounter, 1)
				chi.counter = chi.counter + int(request.ValueToIncrement)
			}

		}(chi)
	}

	select {
	case <-context.Done():
		{
			return nil, fmt.Errorf("Context was done")
		}
	default:
		{
			wg.Wait()
			ok := *okCounter == leng

			if !ok {
				return nil, fmt.Errorf("No consistent")
			}

			s.state.counter = s.state.counter + int(request.ValueToIncrement)

			return &pb.UnaryResponse{
				Ok: true,
			}, nil
		}

	}
}

func (s *server) ParentUnaryCallV2(context context.Context, request *pb.UnaryRequest) (*pb.UnaryResponse, error) {

	s.smaphore <- true
	defer func() { <-s.smaphore }()

	done := make(chan bool)

	leng := int32(len(s.parentData.childs))
	okCounter := new(int32)
	responses := new(int32)

	for _, chi := range s.parentData.childs {

		go func(chi *child, d *chan bool) {
			resp, err := chi.client.UnaryCall(context, request)

			atomic.AddInt32(responses, 1)

			if resp != nil && resp.Ok && err == nil {
				atomic.AddInt32(okCounter, 1)
				chi.counter = chi.counter + int(request.ValueToIncrement)
			}

			if *responses == leng {
				*d <- true
				return
			}

		}(chi, &done)
	}

	<-done

	ok := *okCounter == leng

	if !ok {
		return nil, fmt.Errorf("Inconsistent")
	}

	s.state.counter = s.state.counter + int(request.ValueToIncrement)

	return &pb.UnaryResponse{
		Ok: true,
	}, nil
}

func (s *server) UnaryCall(context context.Context, request *pb.UnaryRequest) (*pb.UnaryResponse, error) {
	s.smaphore <- true
	defer func() { <-s.smaphore }()

	s.state.counter += int(request.ValueToIncrement)

	return &pb.UnaryResponse{
		Ok: true,
	}, nil
}

func NewServerData() *server {

	server := &server{
		globalMutex: sync.Mutex{},
		done:        make(chan bool),
		smaphore:    make(chan bool, 1),
	}

	strs := strings.Split(*ports, ",")

	for _, s := range strs {

		conn, err := grpc.Dial(fmt.Sprintf("localhost:%s", s), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}

		c := &child{
			status:  true,
			client:  pb.NewSerializedServiceClient(conn),
			counter: 0,
		}

		server.parentData.childs = append(server.parentData.childs, c)
	}

	return server

}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *mainPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	server := NewServerData()
	pb.RegisterSerializedServiceServer(s, server)

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
