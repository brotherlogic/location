protoc --proto_path ../../../ -I=./proto --go_out=plugins=grpc:./proto proto/location.proto
mv proto/github.com/brotherlogic/location/proto/* ./proto
