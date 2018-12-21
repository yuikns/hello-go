package examples

import (
	"context"
	"github.com/argcv/stork/log"
	"github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc/grpclog"
	"io/ioutil"
	"os"
	"time"
)

var (
	endpoints = []string{
		"127.0.0.1:2379",
		"127.0.0.1:22379",
		"127.0.0.1:32379",
	}
	dialTimeout    = 5 * time.Second
	requestTimeout = 10 * time.Second
	user           = "root"
	role           = user
	pass           = "1234"
	user2          = "u2"
	role2          = "r2"
	pass2          = "2468"
)

func Ctx() (c context.Context) {
	//c, doCancel := context.WithTimeout(context.TODO(), requestTimeout)
	//doCancel()
	c, _ = context.WithTimeout(context.TODO(), requestTimeout)
	return
}

func InitLogger() {
	gl := grpclog.NewLoggerV2(ioutil.Discard, os.Stderr, os.Stderr)
	clientv3.SetLogger(gl)
}

func Connect(f func(cli *clientv3.Client)) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close() // make sure to close the client
	f(cli)
}
