// +build darwin windows

package main

import (
	"log"
	"strings"
     "github.com/warthog618/gpiod/device/rpi"
	 "fmt"
	"time"
	 )

type Pin struct{
	Name string `storm:"unique"`
	Label string `storm:"unique"`
	ZoneID int `storm:"id"`
	Active bool
	Runtime string
	Pin int
	RuntimeDur time.Duration
}

func (p Pin) On() {

   log.Printf("Pin %s On",p.Name)
}

func (p Pin) Off(){

	log.Printf("Pin %s Off",p.Name)

}

func setupPins(){

	fmt.Println("Setup pins other os")

	var rawpins []Pin

	err := db.All(&rawpins)

	//fmt.Printf("Rawpins %#v\n",rawpins)

	if(err!=nil){}

	for _,par:= range rawpins {


		pinstr:=strings.ToLower(par.Name)

		pin,_:=rpi.Pin(pinstr)
		
		zid := par.ZoneID

		var pa *Pin

		var ok bool

		if pa,ok = pins[zid]; ok {

		}else {
			pa = &Pin{}
		}

		pa.Name = par.Name
		pa.Pin = pin
		pa.Label = par.Label
		pa.ZoneID = par.ZoneID
		pa.Active = par.Active
		pa.Runtime = par.Runtime

		dur,derr := time.ParseDuration(par.Runtime)

		if(derr==nil){

			pa.RuntimeDur = dur

		}

		pins[zid]=pa
	}



}
