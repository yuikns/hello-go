package examples

// https://github.com/coreos/etcd/blob/master/clientv3/example_kv_test.go

import (
	"context"
	"fmt"
	"github.com/argcv/stork/log"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
)

func ExampleKV_put() {
	Connect(func(cli *clientv3.Client) {
		_, err := cli.Put(Ctx(), "sample_key", "sample_value")
		if err != nil {
			log.Fatal(err.Error())
		}
	})
}

func ExampleKV_putErrorHandling() {
	Connect(func(cli *clientv3.Client) {
		_, err := cli.Put(Ctx(), "", "sample_value")
		if err != nil {
			switch err {
			case context.Canceled:
				log.Infof("ctx is canceled by another routine: %v\n", err)
			case context.DeadlineExceeded:
				log.Infof("ctx is attached with a deadline is exceeded: %v\n", err)
			case rpctypes.ErrEmptyKey:
				log.Infof("client-side error: %v\n", err)
			default:
				log.Infof("bad cluster endpoints, which are not etcd servers: %v\n", err)
			}
		}
		// Output: client-side error: etcdserver: key is not provided
	})
}

func ExampleKV_get() {
	Connect(func(cli *clientv3.Client) {
		_, err := cli.Put(context.TODO(), "foo", "bar")
		if err != nil {
			log.Fatal(err)
		}
		resp, err := cli.Get(Ctx(), "foo")
		if err != nil {
			log.Fatal(err)
		}
		for _, ev := range resp.Kvs {
			log.Infof("%s : %s\n", ev.Key, ev.Value)
		}
		// Output: foo : bar
	})
}

func ExampleKV_getWithRev() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	presp, err := cli.Put(context.TODO(), "foo", "bar1")
	if err != nil {
		log.Fatal(err)
	}
	_, err = cli.Put(context.TODO(), "foo", "bar2")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, "foo", clientv3.WithRev(presp.Header.Revision))
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	for _, ev := range resp.Kvs {
		log.Infof("%s : %s\n", ev.Key, ev.Value)
	}
	// Output: foo : bar1
}

func ExampleKV_getSortedPrefix() {
	Connect(func(cli *clientv3.Client) {
		for i := range make([]int, 3) {
			_, err := cli.Put(Ctx(), fmt.Sprintf("key_%d", i), "value")
			if err != nil {
				log.Fatal(err)
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		resp, err := cli.Get(ctx, "key", clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
		cancel()
		if err != nil {
			log.Fatal(err)
		}
		for _, ev := range resp.Kvs {
			log.Infof("%s : %s\n", ev.Key, ev.Value)
		}
		// Output:
		// key_2 : value
		// key_1 : value
		// key_0 : value
	})

}

func ExampleKV_delete() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	// count keys about to be deleted
	gresp, err := cli.Get(ctx, "key", clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}

	// delete the keys
	dresp, err := cli.Delete(ctx, "key", clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted all keys:", int64(len(gresp.Kvs)) == dresp.Deleted)
	// Output:
	// Deleted all keys: true
}

func ExampleKV_compact() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, "foo")
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	compRev := resp.Header.Revision // specify compact revision of your choice

	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	_, err = cli.Compact(ctx, compRev)
	cancel()
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleKV_txn() {
	Connect(func(cli *clientv3.Client) {
		kvc := clientv3.NewKV(cli)
		ctx := Ctx()
		_, err := kvc.Put(ctx, "key", "xyz")
		if err != nil {
			log.Fatal(err)
		}
		_, err = kvc.Txn(ctx).
			// txn value comparisons are lexical
			If(clientv3.Compare(clientv3.Value("key"), ">", "abc")).
			// the "Then" runs, since "xyz" > "abc"
			Then(clientv3.OpPut("key", "XYZ")).
			// the "Else" does not run
			Else(clientv3.OpPut("key", "ABC")).
			Commit()
		if err != nil {
			log.Fatal(err)
		}
		gresp, err := kvc.Get(ctx, "key")
		if err != nil {
			log.Fatal(err)
		}
		for _, ev := range gresp.Kvs {
			log.Infof("%s : %s\n", ev.Key, ev.Value)
		}
		// Output: key : XYZ
	})

}

func ExampleKV_do() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ops := []clientv3.Op{
		clientv3.OpPut("put-key", "123"),
		clientv3.OpGet("put-key"),
		clientv3.OpPut("put-key", "456")}

	for _, op := range ops {
		if _, err := cli.Do(context.TODO(), op); err != nil {
			log.Fatal(err)
		}
	}
}
