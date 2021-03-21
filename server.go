package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
	"html/template"
	_ "embed"
	"context"
)

//go:embed tmpl/index.html
var Index string
//go:embed tmpl/pin.html
var ZoneTmpl string
//go:embed tmpl/sched.html
var SchedTmpl string

type PageData map[string]interface{}

func startWebServer(){

	log.Println("Starting Webserver")

	r := mux.NewRouter()
	r.HandleFunc("/ping", pingHandler)
	r.HandleFunc("/zone/{zone}", zoneHandler)
	r.HandleFunc("/zoneactivate/{zone}",zoneActivateHandler)
	r.HandleFunc("/schedule/{day}", schedHandler)
	r.HandleFunc("/zoneremove/{zone}", zoneRemoveHandler)
	r.HandleFunc("/zonesave", zoneSaveHandler)
	r.HandleFunc("/schedulesave", schedSaveHandler)
	r.HandleFunc("/", indexHandler)
	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":"+viper.GetString("port"),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Listening on port",viper.GetString("port"))
	log.Fatal(srv.ListenAndServe())

}

func pingHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w,"pong")
}

func zoneActivateHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	var pa Pin

	zoneint, _ := strconv.Atoi(vars["zone"])

	err := db.One("ZoneID", zoneint, &pa)

	if (err != nil) {
		fmt.Println(err)
		http.Redirect(w, r, fmt.Sprintf("/?zoneactivate"), 302)
		return
	}

	z := Zone{zoneint,pa.Runtime}

	if(ActiveZone==z){


		log.Println("Zone is already active")
		//already on, turn off!

		ActiveZone.Stop() //turn off the relay
		ActiveZone = Zone{}
		fmt.Fprintf(w,`{"status": "off"}`);
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

	fmt.Fprintf(w,`{"status": "on"}`);

}

func zoneRemoveHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	var pa Pin

	zoneint, _ := strconv.Atoi(vars["zone"])

	err := db.One("ZoneID", zoneint, &pa)

	if (err != nil) {
		fmt.Println(err)
	}

	pa.Active = false

	err = db.Save(&pa)

	var s []Schedule

	db.All(&s)

	for _,sc := range s {
		sc.RemoveZone(zoneint)
    }

	if(err!=nil){fmt.Println(err)}

	setupPins()

	http.Redirect(w, r, fmt.Sprintf("/?zoneremoved"), 302)

}

var decoder  = schema.NewDecoder()

func zoneSaveHandler(w http.ResponseWriter, r *http.Request) {


	r.ParseForm()

	var pa Pin

	err := decoder.Decode(&pa, r.PostForm)

	if err != nil {
		log.Println("Error in POST parameters : ", err)
	}

	pa.Active = true

	fmt.Printf("%#v\n",pa)

	state := db.Save(&pa)




	if(state!=nil){log.Println(state)}

	setupPins()

	http.Redirect(w, r, fmt.Sprintf("/"), 302)

}

func schedSaveHandler(w http.ResponseWriter, r *http.Request) {


	r.ParseForm()

	var s Schedule

	s.Day = r.FormValue("Day")
	s.Starttime = r.FormValue("Starttime")

	frm := make(map[string][]string)

    for k,v := range r.Form {
         frm[k]=v
	}

	runs := frm["Runtime"]

	for zk,z := range frm["ZoneID"] {

		zid,_ := strconv.Atoi(z)
		s.Zones = append(s.Zones,Zone{zid,runs[zk]})

	}


	db.Save(&s)


	http.Redirect(w, r, fmt.Sprintf("/"), 302)

}

func zoneHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	var pa Pin

	zoneint,_ := strconv.Atoi(vars["zone"])

	derr := db.One("ZoneID", zoneint, &pa)

	if(derr!=nil){
		log.Println("Zone Not Found",derr)
	}

	data := PageData{"Title": "Zone","Pin": pa,"Attributes": loadAttributes()}

	var err error
	tmpl := template.New("Zone")

	tmpl,err=tmpl.Parse(ZoneTmpl)

	if(err!=nil) {

		fmt.Println(err)
	}

	err = tmpl.Execute(w, data)

    if(err!=nil){

    	log.Println(err)
	}

}

func schedHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	var sched Schedule

	derr := db.One("Day",vars["day"], &sched)

	if(derr!=nil){
		log.Println("Day Not Found",derr)
	}

	var myzones []int

	for _,p := range pins {


		if(p.Active==true){
			myzones = append(myzones,p.ZoneID)
		}
	}

   sort.Ints(myzones)

	data := PageData{"Title": "Form","MyZones": myzones,"Schedule": sched,"Attributes": loadAttributes()}

	var err error
	tmpl := template.New("Index")
	tmpl,err=tmpl.Parse(SchedTmpl)

	if(err!=nil) {

		fmt.Println(err)
	}

	tmpl.Execute(w, data)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	out := activePins(pins)

	var s []Schedule

	db.All(&s)

	var myzones []int

	for _,p := range pins {


        if(p.Active==true){
		myzones = append(myzones,p.ZoneID)
		}
	}

	log.Printf("Active Schedule %s / Active Zone %s",ActiveSchedule,ActiveZone)

	sort.Ints(myzones)

	sort.Sort(BySchedule(s))

	data := PageData{"Pins": out,"ActiveZone": ActiveZone,"ActiveSchedule": ActiveSchedule,"MyZones": myzones,"Schedule": s}

	tmpl := template.New("Index").Funcs(template.FuncMap{
		"parseHour": parseHour,
	})

    var err error
	tmpl,err=tmpl.Parse(Index)

	if(err!=nil){

		log.Println(err)
	}

	tmpl.Execute(w, data)

}

func activePins(pins map[int]*Pin) (map[int]*Pin){


	out := make(map[int]*Pin)

	for k,v := range pins{

		if(v.Active == true) {

			out[k] = v

		}

	}


	return out


}