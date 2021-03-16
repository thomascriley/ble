module github.com/thomascriley/ble

go 1.13

require (
	github.com/JuulLabs-OSS/cbgo v0.0.2
	github.com/enceve/crypto v0.0.0-20160707101852-34d48bb93815
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/mgutz/logxi v0.0.0-20161027140823-aebf8a7d67ab
	github.com/raff/goble v0.0.0-20200327175727-d63360dcfd80
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae
)

replace github.com/JuulLabs-OSS/cbgo => github.com/thomascriley/cbgo v0.0.3-0.20210309070341-a5fcee8c38af
