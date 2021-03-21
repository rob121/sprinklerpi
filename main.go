package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/asdine/storm/v3"
	"log"
	"strconv"
	"time"
)

var pins map[int]*Pin
var db *storm.DB

func main(){

	pins = make(map[int]*Pin)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
    setupConfig()
	initDb()
	setupPins()
	go setupHomekit()
	go startSchedule()
	startWebServer()
	select{}

}




func setupConfig(){

	log.Println("Setting Configuration")

	viper.SetDefault("AppName","sprinklerpi")
	viper.SetDefault("HttpActionTimeout","15s")
	viper.SetDefault("GpioPinCt",40)
	viper.SetDefault("Chip","gpiochip0")
	viper.SetDefault("Debug",false)
	viper.SetDefault("port","8000")
	viper.SetConfigType("json")
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/sprinklerpi/")   // path to look for the config file in
	viper.AddConfigPath("$HOME/.sprinklerpi")  // call multiple times to add many search paths
	            // optionally look for config in the working directory
	viper.WatchConfig()


	viper.OnConfigChange(func(e fsnotify.Event) {

		log.Println("Config file changed:", e.Name)
	})

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			panic(err)
		}
	}

	log.Println(viper.ConfigFileUsed())
}

type Attributes struct {
	Zones []int
	Pins []string
	Days []string
}

func loadAttributes() Attributes {

	var defpins []string

	for i:=1;i<=viper.GetInt("GpioPinCt");i++ {
		defpins = append(defpins,"GPIO"+strconv.Itoa(i))
	}



	days := []string{"Sunday","Monday","Tuesday","Wednesday","Thursday","Friday","Saturday"}

	return Attributes{Days: days,Zones: []int{1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20}, Pins: defpins}

}

func parseHour(hour string) string{
	layout2 := "03:04 PM"
	layout1 := "15:04"
	t, err := time.Parse(layout1, hour)
	if err != nil {
		log.Println(err)
        return ""
	}

	return t.Format(layout2)


}

var daysOfWeek = map[string]time.Weekday{
	"Sunday":    time.Sunday,
	"Monday":    time.Monday,
	"Tuesday":   time.Tuesday,
	"Wednesday": time.Wednesday,
	"Thursday":  time.Thursday,
	"Friday":    time.Friday,
	"Saturday":  time.Saturday,
}

func parseWeekday(v string) (time.Weekday, error) {
	if d, ok := daysOfWeek[v]; ok {
		return d, nil
	}

	return time.Sunday, fmt.Errorf("invalid weekday '%s'", v)
}
