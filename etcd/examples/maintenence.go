package examples

// https://github.com/coreos/etcd/blob/master/clientv3/example_maintenence_test.go
import (
	"github.com/argcv/stork/log"
	"github.com/coreos/etcd/clientv3"
)

func ExampleMaintenance_status() {
	Connect(func(cli *clientv3.Client) {
		for _, ep := range endpoints {
			resp, err := cli.Status(Ctx(), ep)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			log.Infof("endpoint: %s / Leader: %v\n", ep, resp.Header.MemberId == resp.Leader)
		}
	})
	// endpoint: localhost:2379 / Leader: false
	// endpoint: localhost:22379 / Leader: false
	// endpoint: localhost:32379 / Leader: true
}

func ExampleMaintenance_defragment() {
	Connect(func(cli *clientv3.Client) {
		for _, ep := range endpoints {
			if resp, err := cli.Defragment(Ctx(), ep); err != nil {
				log.Error(err.Error())
			} else {
				log.Infof("%v", resp)
			}
		}
	})
}

func ExampleMaintenance() {
	ExampleMaintenance_status()
	ExampleMaintenance_defragment()
}
