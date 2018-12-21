package main

import (
	"github.com/yuikns/hello-go/etcd/examples"
)

func main() {
	examples.InitLogger()
	examples.ExampleBase()
	examples.ExampleAuth()
	examples.ExampleCluster_memberList()
	examples.ExampleMaintenance()
}
