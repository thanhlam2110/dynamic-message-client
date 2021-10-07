package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"time"

	"grpc-client/proto"
	"log"

	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	address = "localhost:8000"
)

func main() {
	// testPutData()
	// testGetData()
	// testPutDataStream()
	testGetDataStream()
}

func testPutData() {
	putData("int", []byte("1"))
	putData("float", []byte("1.1"))
	putData("bool", []byte("true"))
	putData("string", []byte("this is a string msg"))
	putData("json", []byte("{\"msg\":\"this is a json msg\", \"other\":1}"))
}

func testPutDataStream() {
	var testData = map[string][]byte{
		"int":    []byte("1"),
		"float":  []byte("1.1"),
		"bool":   []byte("true"),
		"string": []byte("this is a string msg"),
	}
	putDataStream(testData)
}

func testGetData() {
	getData("int", []byte("1"))
}

func testGetDataStream() {
	var testData = map[string][]byte{
		"int":    []byte("1"),
		"float":  []byte("1.1"),
		"bool":   []byte("true"),
		"string": []byte("this is a string msg"),
	}
	getDataStream(testData)
}

func putData(valueType string, value []byte) error {

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewDataServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	reqData := &proto.PutDataRequest{
		Data: &any.Any{Value: value},
		Type: valueType,
	}

	resp, err := client.PutData(ctx, reqData)
	if err != nil {
		log.Printf("put data error. %v", err)
		return err
	}

	if resp.Err != 0 {
		log.Printf("put data exec error. %v", resp.Desc)
		return err
	}

	return nil
}

// values : key is type, value is value
func putDataStream(values map[string][]byte) error {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewDataServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	stream, err := client.PutDataStream(ctx)
	if err != nil {
		log.Printf("put data stream get stream err :%v", err)
		return err
	}

	for valueType, value := range values {
		reqData := &proto.PutDataStreamRequest{
			Data: &any.Any{Value: value},
			Type: valueType,
		}

		err := stream.Send(reqData)
		if err != nil {
			log.Printf("put data stream send error. %v", err)
			return err
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("put data stream close error. %v", err)
		return err
	}

	if resp.Err != 0 {
		log.Printf("put data stream exec error. %v", resp.Desc)
		return err
	}

	return nil
}

func getData(valueType string, value []byte) error {

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewDataServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	reqData := &proto.GetDataRequest{
		Data: &any.Any{Value: value},
		Type: valueType,
	}

	stream, err := client.GetData(ctx, reqData)
	if err != nil {
		log.Printf("get data error. %v", err)
		return err
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("get data recv data error: %v", err)
			return err
		}

		recvData, err := parseAnyData(msg.Type, msg.Data)
		if err != nil {
			log.Printf("get data parse data error: %v", err)
			continue
		}
		log.Printf("recv from server: %v", recvData)

	}

	return nil
}

// values : key is type, value is value
// send then recv
// send and recv are extracted and put into the coroutine to make them independent
func getDataStream(values map[string][]byte) error {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewDataServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	stream, err := client.GetDataStream(ctx)
	if err != nil {
		log.Printf("get data stream get stream err :%v", err)
		return err
	}
	waitc := make(chan struct{})
	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Printf("get data stream recv data error: %v", err)
				break
			}
			recvData, err := parseAnyData(msg.Type, msg.Data)
			if err != nil {
				log.Printf("get data stream parse data error: %v", err)
				continue
			}
			log.Printf("recv from server: %v", recvData)
		}
	}()

	for valueType, value := range values {

		reqData := &proto.GetDataStreamRequest{
			Data: &any.Any{Value: value},
			Type: valueType,
		}

		err := stream.Send(reqData)
		if err != nil {
			log.Printf("get data stream send error. %v", err)
			return err
		}
	}
	stream.CloseSend()
	<-waitc
	return nil
}

func parseAnyData(valueType string, anyData *anypb.Any) (interface{}, error) {

	if anyData == nil {
		return nil, errors.New("illegal param")
	}

	var err error
	var reqData interface{}
	switch valueType {
	case "int":
		reqData, err = strconv.ParseInt(string(anyData.Value), 10, 64)
	case "float":
		reqData, err = strconv.ParseFloat(string(anyData.Value), 10)
	case "json":
		err = json.Unmarshal(anyData.Value, &reqData)
	case "bool":
		reqData, err = strconv.ParseBool(string(anyData.Value))
	case "string":
		reqData = string(anyData.Value)
	default:
		reqData = string(anyData.Value)
	}

	return reqData, err
}
