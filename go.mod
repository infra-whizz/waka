module github.com/infra-whizz/waka

go 1.14

require (
	github.com/antonfisher/nested-logrus-formatter v1.1.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/infra-whizz/wzlib v0.0.0-20200630192324-59872deadc1a
	github.com/isbm/go-nanoconf v0.0.0-20200623180822-caf90de1965e
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/urfave/cli/v2 v2.2.0
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
)

replace github.com/infra-whizz/wzlib => ../wzlib
