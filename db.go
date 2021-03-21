package main


import(
	"github.com/kirsle/configdir"
 	"github.com/asdine/storm/v3"
	"path/filepath"
	"fmt"

)

var configPath string


func initDb(){

	var err error
	configPath = configdir.LocalConfig("sprinklerpi")

	err = configdir.MakePath(configPath) // Ensure it exists.

	if err != nil {
	panic(err)
	}

	fmt.Println(configPath)

	db, err = storm.Open(filepath.Join(configPath, "sprinklerpi.db"))

	if err!=nil {

	panic(err)

	}

}
