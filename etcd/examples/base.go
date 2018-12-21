package examples

// https://github.com/coreos/etcd/blob/master/clientv3/example_test.go
import (
	"github.com/argcv/stork/log"
	"github.com/coreos/etcd/clientv3"
	"github.com/davecgh/go-spew/spew"
)

func ExampleBase() {
	Connect(func(cli *clientv3.Client) {
		var err error
		_, err = cli.Put(Ctx(), "foo", "bar")
		if err != nil {
			log.Fatal(err)
		}
		_, err = cli.Put(Ctx(), "foo", "bar2")
		if err != nil {
			log.Fatal(err)
		}
		var gr *clientv3.GetResponse
		gr, err = cli.Get(Ctx(), "foo")
		if err != nil {
			log.Info(err)
		} else {
			log.Info(spew.Sdump(gr))
			for i, kv := range gr.Kvs {
				log.Infof("#[%v] cv: %v, mv: %v, ver: %v",
					i, kv.CreateRevision, kv.ModRevision, kv.Version)
				log.Infof("#[%v] k: %v, v: %v, l: %v",
					i, string(kv.Key[:]), string(kv.Value[:]), kv.Lease)
				log.Infof("#[%v] str: %v", i, kv.String())
			}
			log.Info(gr.Kvs)
		}
	})

}
