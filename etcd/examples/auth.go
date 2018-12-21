package examples

// https://github.com/coreos/etcd/blob/master/clientv3/example_auth_test.go

import (
	"github.com/argcv/stork/log"
	"github.com/coreos/etcd/clientv3"
)

func ExampleAuth() {
	Connect(func(cli *clientv3.Client) {
		var err error
		if _, err = cli.RoleAdd(Ctx(), role); err != nil {
			log.Fatal(err)
		}
		if _, err = cli.UserAdd(Ctx(), user, pass); err != nil {
			log.Fatal(err)
		}
		if _, err = cli.UserGrantRole(Ctx(), user, role); err != nil {
			log.Fatal(err)
		}

		if _, err = cli.RoleAdd(Ctx(), role2); err != nil {
			log.Fatal(err)
		}

		if _, err = cli.RoleGrantPermission(
			Ctx(),
			role2, // role name
			"foo", // key
			"zoo", // range end
			clientv3.PermissionType(clientv3.PermReadWrite),
		); err != nil {
			log.Fatal(err)
		}
		if _, err = cli.UserAdd(Ctx(), user2, pass2); err != nil {
			log.Fatal(err)
		}
		if _, err = cli.UserGrantRole(Ctx(), user2, role2); err != nil {
			log.Fatal(err)
		}
		if _, err = cli.AuthEnable(Ctx()); err != nil {
			log.Fatal(err)
		}

		cliAuth, err := clientv3.New(clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: dialTimeout,
			Username:    user2,
			Password:    pass2,
		})
		if err != nil {
			log.Fatal(err)
		}
		defer cliAuth.Close()

		if _, err = cliAuth.Put(Ctx(), "foo1", "bar"); err != nil {
			log.Fatal(err)
		}

		_, err = cliAuth.Txn(Ctx()).
			If(clientv3.Compare(clientv3.Value("zoo1"), ">", "abc")).
			Then(clientv3.OpPut("zoo1", "XYZ")).
			Else(clientv3.OpPut("zoo1", "ABC")).
			Commit()
		log.Info("[Expected]", err)

		// now check the permission with the root account
		rootCli, err := clientv3.New(clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: dialTimeout,
			Username:    user,
			Password:    pass,
		})
		if err != nil {
			log.Fatal(err)
		}
		defer rootCli.Close()

		resp, err := rootCli.RoleGet(Ctx(), role2)
		if err != nil {
			log.Fatal(err)
		}
		log.Infof("[Expected] user u permission: key %q, range end %q\n", resp.Perm[0].Key, resp.Perm[0].RangeEnd)

		if _, err = rootCli.AuthDisable(Ctx()); err != nil {
			log.Fatal(err)
		}
		// Output: etcdserver: permission denied
		// user u permission: key "foo", range end "zoo"

		cli.RoleDelete(Ctx(), role)
		if _, err := cli.RoleDelete(Ctx(), role); err != nil {
			log.Info("[Expected]", err)
		}
		if _, err := cli.RoleDelete(Ctx(), role2); err != nil {
			log.Fatal(err)
		}
		if _, err := cli.UserDelete(Ctx(), user); err != nil {
			log.Fatal(err)
		}
		if _, err := cli.UserDelete(Ctx(), user2); err != nil {
			log.Fatal(err)
		}
	})

}
