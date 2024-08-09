module github.com/atmatm9182/silly-hippos/server

go 1.22.5

replace github.com/atmatm9182/silly-hippos/common => ../common/

require (
	github.com/atmatm9182/silly-hippos/common v0.0.0-20240805174613-b77575719b11
	google.golang.org/protobuf v1.34.2
)

require github.com/atmatm9182/silly-hippos/hwg v0.0.0-20240805195008-2c62a2122f9a // indirect

replace github.com/atmatm9182/silly-hippos/hwg => ../hwg
