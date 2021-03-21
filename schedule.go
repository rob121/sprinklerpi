package main

import (
	"context"
	"log"
	"fmt"
	"time"
)
var ZoneChan chan ZoneState
var ActiveZone Zone
var ActiveSchedule Schedule


func init(){

	ZoneChan = make(chan ZoneState)

}

type ZoneState struct{
	Zone Zone
	State int
}

type Zone struct {
	ZoneID int
	Runtime string
}

type Schedule struct {
	Starttime string
	Day string `storm:"id"`
	Zones []Zone
	Running bool
}


type BySchedule []Schedule

func (s BySchedule) Len() int {
	return len(s)
}

func (s BySchedule) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s  BySchedule) Less(i, j int) bool {

	d1,_ := parseWeekday(s[i].Day)
    d2,_ := parseWeekday(s[j].Day)
	return d1<d2
}

func (s Schedule) String() string{


	return fmt.Sprintf("%s @ %s",s.Day,s.Starttime)

}

func (z Zone) String() string{


	return fmt.Sprintf("%d for %s",z.ZoneID,z.Runtime)

}

func (z Zone) Stop(){

	for _,p := range pins {

		if(p.ZoneID == z.ZoneID){

			//start
			log.Printf("Stop Zone:%d",p.ZoneID)
			p.Off()
			ZoneChan <- ZoneState{Zone{ZoneID:p.ZoneID},0}
		}

	}

	ActiveZone = Zone{}

}

func (z Zone) Run(ctx context.Context){
 //get the pin and start it (and stop the others)

	ActiveZone = z

	for _,p := range pins {

		if(p.ZoneID != z.ZoneID){

			//stop
            p.Off()
			ZoneChan <- ZoneState{Zone{ZoneID:p.ZoneID},0}
			log.Printf("Stopping Zone:%d",p.ZoneID)
		}

	}

	for _,p := range pins {

		if(p.ZoneID == z.ZoneID){

			//start
			p.On()
			ZoneChan <- ZoneState{Zone{ZoneID:p.ZoneID},1}
			log.Printf("Starting Zone:%d",p.ZoneID)
		}

	}

	select {
	case <-ctx.Done():
      log.Printf("Zone Run %d Stopped...\n",z.ZoneID)
      z.Stop()
	}



}
func (s Schedule) RemoveZone(id int){

  var zo []Zone

  for _,z := range s.Zones {

  	 if(z.ZoneID != id){

  	 	zo = append(zo,z)

	 }



  }

  s.Zones = zo

  db.Save(&s)

}



func (s Schedule) Start(){

	s.Running = true

	for _,z := range s.Zones {


		ctx, cancel := context.WithCancel(context.Background())

		go z.Run(ctx)

		dur,err := time.ParseDuration(z.Runtime)


		if(err!=nil){

			log.Println("Invalid time, going to next zone...")
			continue
		}

		log.Printf("Running Zone %d until %s using duration of %s",z.ZoneID,time.Now().Add(dur).Format("03:04 PM"),z.Runtime)

		select {
		case <-time.After(dur):
			log.Println("Zone Complete, moving to next zone...")
			cancel()
		case <-time.After(2 * time.Hour):
			cancel()
			log.Println("Zone time out")
		}

		//we either timed out or ran the duration, next zone!

	}

	s.Running = false
	log.Printf("Schedule %s complete\n",s)
	//TODO stop all zones


}

func startSchedule() {

	ticker := time.NewTicker(1*time.Minute)

	checkSchedule()

	for range ticker.C {

		//find an active schedule

        checkSchedule()

	}

}



func checkSchedule(){

	var sched []Schedule

	db.All(&sched)

	//log.Printf("Checking %s at %s\n",time.Now().Weekday(),time.Now().Format("15:04"))

	for k,s := range sched {


		if(isScheduleTime(&s)){



			ActiveSchedule = s
			go s.Start()

			log.Printf("Starting Schedule %s\n",s)

		}else{

			sched[k].Running = false
		}

	}

}

func isScheduleTime(s *Schedule) (bool){
	//log.Printf("Checking Schedule match %s at %s = %s at %s\n",s.Day,s.Starttime,time.Now().Weekday(),time.Now().Format("15:04"))

	if( s.Day == time.Now().Weekday().String()) {

		//check if it's time

		if( string(time.Now().Format("15:04"))  == s.Starttime ) {

			 return true

		}else{

			//log.Println("time not match")
		}

	}else{

		//log.Println("day not match")
	}

	return false

}