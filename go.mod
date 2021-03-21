module github.com/rob121/sprinklerpi

go 1.16

replace github.com/brutella/hc => /Users/ralfonso/go/src/github.com/rob121/sprinklerpi/hc

require (
	github.com/asdine/storm/v3 v3.2.1
	github.com/brutella/hc v1.2.4
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/schema v1.2.0
	github.com/kirsle/configdir v0.0.0-20170128060238-e45d2f54772f
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.7.1
	github.com/warthog618/gpiod v0.6.0
)
