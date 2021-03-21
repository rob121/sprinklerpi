package main

import (
	"github.com/brutella/hc"
	//llog "github.com/brutella/hc/log"
	"github.com/brutella/hc/accessory"
	"time"
"context"
	//	"github.com/brutella/hc/service"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func setupHomekit() {

   // llog.Debug.Enable()


	log.Info("Setting Up Homekit Device: ", fmt.Sprintf("Sprinkler System"))

	info := accessory.Info{Model: "V1", ID: uint64(1), Manufacturer: "FonziCorp", Name: "Sprinkler System"}

	acc := accessory.NewSprinklerSystem(info)

	//bridge := accessory.NewBridge(accessory.Info{Name: "Sprinkler Bridge", ID: 1})

	config := hc.Config{Pin: "12344321", Port: "12345", StoragePath: filepath.Join(configPath, "homekit.db")}
	t, err := hc.NewIPTransport(config, acc.Accessory)
	if err != nil {
		log.Error(err)
	}

	acc.Zone1.Active.OnValueRemoteUpdate(func(active int) { homekitValve(1, active) })
	acc.Zone2.Active.OnValueRemoteUpdate(func(active int) { homekitValve(2, active) })
	acc.Zone3.Active.OnValueRemoteUpdate(func(active int) { homekitValve(3, active) })
	acc.Zone4.Active.OnValueRemoteUpdate(func(active int) { homekitValve(4, active) })
	acc.Zone5.Active.OnValueRemoteUpdate(func(active int) { homekitValve(5, active) })
	acc.Zone6.Active.OnValueRemoteUpdate(func(active int) { homekitValve(6, active) })

	go func() {

		for z := range ZoneChan {

			log.Infof("Got Zone %d state %d",z.Zone.ZoneID,z.State)

			switch z.Zone.ZoneID {

			case 1:
				acc.Zone1.Active.SetValue(z.State)
			case 2:
				acc.Zone2.Active.SetValue(z.State)
			case 3:
				acc.Zone3.Active.SetValue(z.State)
			case 4:
				acc.Zone4.Active.SetValue(z.State)
			case 5:
				acc.Zone5.Active.SetValue(z.State)
			case 6:
				acc.Zone6.Active.SetValue(z.State)

			}

		}
	}()

	hc.OnTermination(func() {
		log.Info("Interrupt Received")
		os.Exit(1)
		<-t.Stop()
	})

	go t.Start()
}

func homekitValve(zoneint int, active int) {

	log.Infof("Zone %d state %d", zoneint, active)

	var pa Pin

	err := db.One("ZoneID", zoneint, &pa)

	if (err != nil) {
		fmt.Println(err)
		return
	}

	z := Zone{zoneint,pa.Runtime}

	if(active==0){

 		z.Stop()
        return
	}


	go func() {

		ctx, cancel := context.WithCancel(context.Background())
		go z.Run(ctx)

		dur, err := time.ParseDuration(pa.Runtime)

		if (err != nil) {

			log.Println("Invalid time, going to next zone...")
			return
		}

		log.Printf("Running Zone %d until %s using duration of %s", z.ZoneID, time.Now().Add(dur).Format("03:04 PM"), pa.Runtime)

		select {
		case <-time.After(dur):
			log.Println("Zone Complete...")
			cancel()
		case <-time.After(2 * time.Hour):
			cancel()
			log.Println("Zone time out")
		}

	}()


}
