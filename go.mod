module github.com/infra-whizz/waka

go 1.14

require (
	github.com/antonfisher/nested-logrus-formatter v1.3.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/infra-whizz/wzlib v0.0.0-20210302184443-9c300336f5e1
	github.com/infra-whizz/wzmodlib v0.0.0-20200629192210-dc7f33b52b2b // indirect
	github.com/isbm/go-nanoconf v0.0.0-20200623180822-caf90de1965e
	github.com/isbm/go-shutil v0.0.0-20200707163617-60e3684d72ba
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/magefile/mage v1.11.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sirupsen/logrus v1.8.0 // indirect
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/sys v0.0.0-20210305230114-8fe3ee5dd75b // indirect
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/infra-whizz/wzlib => ../wzlib

replace github.com/isbm/go-shutil => ../../go-shutil
