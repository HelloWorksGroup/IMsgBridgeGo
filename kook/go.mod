module kookNode

go 1.20

require (
	github.com/HelloWorksGroup/IMSuperGroup/imnode v0.0.0-00010101000000-000000000000
	github.com/lonelyevil/kook v0.0.33
	github.com/lonelyevil/kook/log_adapter/plog v0.0.31
	github.com/phuslu/log v1.0.83
)

require (
	github.com/bits-and-blooms/bitset v1.2.2 // indirect
	github.com/bits-and-blooms/bloom/v3 v3.2.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
)


replace github.com/HelloWorksGroup/IMSuperGroup/imnode => ../imnode
