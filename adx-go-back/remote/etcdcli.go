package remote

import (
	"adx-go/common"
	"context"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
)

var etcdserver = common.GetEnvDef("ETCD_ENDPOINTS", "192.168.1.214:2379")

func GetVal(key string) string {
	var err error
	var cli *clientv3.Client
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdserver},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	var result string
	if len(key) < 1 {
		return result
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := cli.Get(ctx, key)
	if err != nil {
		log.Println(err)
	}
	for _, ev := range resp.Kvs {
		log.Printf("%s : %s\n", ev.Key, ev.Value)
		result = string(ev.Value)
	}
	return result
}
