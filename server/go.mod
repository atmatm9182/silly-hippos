module github.com/atmatm9182/silly-hippos/server

go 1.22.5

replace github.com/atmatm9182/silly-hippos/common => ../common/

require (
	github.com/atmatm9182/silly-hippos/common v0.0.0-00010101000000-000000000000
	google.golang.org/protobuf v1.34.2
)
