package examples

// https://github.com/coreos/etcd/blob/master/clientv3/example_cluster_test.go
import (
	"fmt"
	"github.com/argcv/stork/log"
	"github.com/coreos/etcd/clientv3"
)

func ExampleCluster_memberList() {
	Connect(func(cli *clientv3.Client) {
		resp, err := cli.MemberList(Ctx())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("members:", len(resp.Members))
		// Output: members: 3
		for i, mem := range resp.Members {
			log.Info(i, mem)
		}
	})
}

func ExampleCluster_memberAdd() {
	Connect(func(cli *clientv3.Client) {
		peerURLs := endpoints[2:]
		mresp, err := cli.MemberAdd(Ctx(), peerURLs)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("added member.PeerURLs:", mresp.Member.PeerURLs)
		// added member.PeerURLs: [http://localhost:32380]
	})

}

func ExampleCluster_memberRemove() {
	Connect(func(cli *clientv3.Client) {
		resp, err := cli.MemberList(Ctx())
		if err != nil {
			log.Fatal(err)
		}

		_, err = cli.MemberRemove(Ctx(), resp.Members[0].ID)
		if err != nil {
			log.Fatal(err)
		}
	})

}

func ExampleCluster_memberUpdate() {
	Connect(func(cli *clientv3.Client) {
		resp, err := cli.MemberList(Ctx())
		if err != nil {
			log.Fatal(err)
		}

		peerURLs := []string{"http://localhost:12380"}
		_, err = cli.MemberUpdate(Ctx(), resp.Members[0].ID, peerURLs)
		if err != nil {
			log.Fatal(err)
		}
	})
}
