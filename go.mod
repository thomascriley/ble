module github.com/thomascriley/ble

go 1.18

require (
	github.com/JuulLabs-OSS/cbgo v0.0.2
	github.com/raff/goble v0.0.0-20200327175727-d63360dcfd80
	golang.org/x/sys v0.20.0
)

require github.com/sirupsen/logrus v1.9.3 // indirect

replace github.com/JuulLabs-OSS/cbgo v0.0.2 => github.com/thomascriley/cbgo v0.0.4
