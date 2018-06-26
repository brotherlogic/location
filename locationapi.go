package main

import (
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/location/proto"
)

// AddLocation adds the user location
func (s *Server) AddLocation(ctx context.Context, req *pb.AddLocationRequest) (*pb.AddLocationResponse, error) {
	req.GetLocation().Time = time.Now().Unix()
	s.config.Locations = append(s.config.Locations, req.GetLocation())
	s.save()

	return &pb.AddLocationResponse{}, nil
}

// GetLocation gets the most recent user location
func (s *Server) GetLocation(ctx context.Context, req *pb.GetLocationRequest) (*pb.GetLocationResponse, error) {
	var bestLocation *pb.Location
	bestTime := int64(0)
	for _, l := range s.config.Locations {
		if l.Name == req.Name && l.Time > bestTime {
			bestLocation = l
			bestTime = l.Time
		}
	}

	return &pb.GetLocationResponse{Location: bestLocation}, nil
}
