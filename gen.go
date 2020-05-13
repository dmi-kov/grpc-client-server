package main

/////////////////////////////////////////////////
// This file used for auto-generation purposes //
/////////////////////////////////////////////////

//go:generate sh -c "protoc -I api/ -I${GOPATH}/src --go_out=plugins=grpc:api api/api.proto"
