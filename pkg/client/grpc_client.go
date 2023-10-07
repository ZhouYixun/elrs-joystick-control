package client

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/kaack/elrs-joystick-control/pkg/proto/generated/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
	"os"
	"time"
)

func Init(txServerPortName, configFilePath string, txServerPortBaudRate int) {
	if (len(txServerPortName) != 0 || len(configFilePath) != 0) && txServerPortBaudRate != 0 {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var err error
		var conn *grpc.ClientConn
		if conn, err = grpc.DialContext(ctx, "localhost:10000"); err != nil {
			panic(err)
		}

		client := pb.NewJoystickControlClient(conn)
		var res *pb.Empty

		if len(txServerPortName) != 0 && txServerPortBaudRate != 0 {
			if res, err = client.StartLink(ctx, &pb.StartLinkReq{
				Port:     txServerPortName,
				BaudRate: int32(txServerPortBaudRate),
			}); err != nil {
				panic(err)
			}

			fmt.Printf("%v", res)
		}

		if len(configFilePath) != 0 {

			var configJson []byte
			configJson, err = os.ReadFile(configFilePath)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(configJson))

			var configPb structpb.Struct
			m := jsonpb.Unmarshaler{}
			if err = m.Unmarshal(bytes.NewReader(configJson), &configPb); err != nil {
				panic(err)
			}

			if res, err = client.SetConfig(ctx, &pb.SetConfigReq{
				Config: &configPb,
			}); err != nil {
				panic(err)
			}

			fmt.Printf("%v", res)
		}
	}
}
