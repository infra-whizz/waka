module github.com/infra-whizz/waka

go 1.14

require (
	github.com/antonfisher/nested-logrus-formatter v1.3.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/infra-whizz/wzlib v0.0.0-20210306212611-2af49aea1704
	github.com/isbm/go-nanoconf v0.0.0-20200623180822-caf90de1965e
	github.com/isbm/go-shutil v0.0.0-20200707163617-60e3684d72ba
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/sys v0.0.0-20210326220804-49726bf1d181 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

// replace github.com/infra-whizz/wzlib => ../wzlib
// replace github.com/isbm/go-shutil => ../../go-shutil
