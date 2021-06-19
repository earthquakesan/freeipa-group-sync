module github.com/earthquakesan/freeipa-group-sync

go 1.16

require (
	github.com/ubccr/goipa v0.0.5 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

// For development point goipa to my local fork
replace github.com/ubccr/goipa v0.0.5 => ../goipa
