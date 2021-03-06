package etcdv3

import (
	"github.com/coder2z/g-saber/xcfg"
	"github.com/coder2z/g-saber/xflag"
	"net/url"
	"time"

	"go.etcd.io/etcd/clientv3"
)

// DataSourceEtcd defines etcd scheme
const DataSourceEtcd = "etcd"

type etcd struct{}

func New() *etcd {
	return new(etcd)
}

func (e etcd) Register() (string, func() xcfg.DataSource) {
	return DataSourceEtcd, func() xcfg.DataSource {
		var (
			configAddr = xflag.String("xcfg")
			watch      = xflag.Bool("watch")
		)
		if configAddr == "" {
			return nil
		}
		// configAddr is a string in this format:
		// etcd://ip:port?username=XXX&password=XXX&key=key

		urlObj, err := url.Parse(configAddr)
		if err != nil {
			return nil
		}
		etcdConf := clientv3.Config{
			DialKeepAliveTime:    10 * time.Second,
			DialKeepAliveTimeout: 3 * time.Second,
		}
		etcdConf.Endpoints = []string{urlObj.Host}
		etcdConf.Username = urlObj.Query().Get("username")
		etcdConf.Password = urlObj.Query().Get("password")
		client, err := clientv3.New(etcdConf)
		if err != nil {
			return nil
		}
		return NewDataSource(client, urlObj.Query().Get("key"), watch)
	}
}
