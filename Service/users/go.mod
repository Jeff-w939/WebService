module users

go 1.18

require (
	github.com/aliyun/alibaba-cloud-sdk-go v0.0.0-20190808125512-07798873deee
	github.com/golang/protobuf v1.5.2
	github.com/gomodule/redigo/redis v0.0.1
	//github.com/gomodule/redigo/redis v0.0.1
	github.com/jinzhu/gorm v1.9.16
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins/registry/consul v0.0.0-20200119172437-4fe21aa238fd
	github.com/micro/micro/v3 v3.11.0
	google.golang.org/protobuf v1.28.0
)

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
