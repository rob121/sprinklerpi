// +build linux

package main

import(
	"github.com/warthog618/gpiod"
	"fmt"
	"time"
	"log"
	"strings"
	"github.com/warthog618/gpiod/device/rpi"
	)




type Pin struct{
	Name string `storm:"unique"`
	Label string `storm:"unique"`
	ZoneID int `storm:"id"`
	Active bool
	PinLine    *gpiod.Line `json:"-"`
	Runtime string
	Pin int
	RuntimeDur time.Duration
}

func (p Pin) On() {

	log.Printf("Pin %s On",p.Name)
	p.PinLine.SetValue(0)

}

func (p Pin) Off(){

	p.PinLine.SetValue(1)
	log.Printf("Pin %s Off",p.Name)
}

var c *gpiod.Chip

func init(){


	var err error
	c, err = gpiod.NewChip("gpiochip0")

    if(err!=nil){

    	panic(err)

	}

}

func setupPins(){

	fmt.Println("Setup pins linux")

	//l.SetValue(0) //on
	//l.SetValue(1) //off


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

			l, err := c.RequestLine(pin, gpiod.AsOutput(1))
			pa.PinLine = l
			if(err!=nil){

				log.Println(par.Name,err)
			}
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


