package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/goserver/utils"
	"github.com/brotherlogic/keystore/client"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pbg "github.com/brotherlogic/goserver/proto"
	pbled "github.com/brotherlogic/led/proto"
	pb "github.com/brotherlogic/location/proto"
)

const (
	// KEY - where the locations are stored
	KEY = "/github.com/brotherlogic/location/config"
)

//Server main server type
type Server struct {
	*goserver.GoServer
	config *pb.Config
	writer writer
	counts int64
}

// Init builds the server
func Init() *Server {
	s := &Server{
		&goserver.GoServer{},
		&pb.Config{},
		&prodWriter{},
		0,
	}
	s.writer = &prodWriter{Log: s.Log}
	return s
}

type writer interface {
	writeToLed(ctx context.Context, top, bot string)
}

type prodWriter struct {
	Log func(text string)
}

func (p *prodWriter) writeToLed(ctx context.Context, top, bot string) {
	ip, port, err := utils.Resolve("led")
	if err == nil {
		conn, err := grpc.Dial(ip+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
		if err == nil {
			client := pbled.NewLedServiceClient(conn)
			r, err := client.Display(ctx, &pbled.DisplayRequest{TopLine: top, BottomLine: bot})
			p.Log(fmt.Sprintf("Written %v and %v", r, err))
		}
	}
}

func (s *Server) save(ctx context.Context) {
	s.KSclient.Save(ctx, KEY, s.config)
}

func (s *Server) load(ctx context.Context) error {
	config := &pb.Config{}
	data, _, err := s.KSclient.Read(ctx, KEY, config)

	if err != nil {
		return err
	}

	s.config = data.(*pb.Config)
	return nil
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	pb.RegisterLocationServiceServer(server, s)
}

// ReportHealth alerts if we're not healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Mote promotes/demotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	if master {
		err := s.load(ctx)
		return err
	}

	return nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	return []*pbg.State{
		&pbg.State{Key: "counts", Value: s.counts},
	}
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	flag.Parse()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	server := Init()
	server.GoServer.KSclient = *keystoreclient.GetClient(server.GetIP)
	server.PrepServer()
	server.Register = server

	server.RegisterServer("location", false)
	server.Log("Starting!")
	server.Serve()
}
