module xstorage

go 1.16

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/gorilla/websocket v1.4.2
	github.com/json-iterator/go v1.1.11
	github.com/mitchellh/cli v1.1.0
	github.com/smallnest/rpcx v0.0.0-20210120041900-c2830baacdb1
	github.com/wlgd/xproto v0.0.0-00010101000000-000000000000
)

replace github.com/wlgd/xproto => ../xproto
