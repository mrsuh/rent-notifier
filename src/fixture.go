package main

import (
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"os"
	"log"
	"rent-notifier/src/db"
	"github.com/mrsuh/cli-config"
)

func connectDb() *dbal.DBAL {
	conf_instance := config.GetInstance()

	err := conf_instance.Init()

	if err != nil {
		log.Fatal(err)
	}

	conf := conf_instance.Get()

	connection := dbal.NewConnection(conf["database.dsn"].(string))
	defer connection.Session.Close()

	return &dbal.DBAL{DB: connection.Session.DB(connection.Database)}
}

func loadCities(db *dbal.DBAL, path string) {
	data, read_err := ioutil.ReadFile(path)

	if read_err != nil {
		fmt.Println(read_err)
		os.Exit(1)
	}

	json := make([]map[string]interface{}, 0)

	parse_err := yaml.Unmarshal(data, &json)
	if parse_err != nil {
		fmt.Println(parse_err)
		os.Exit(1)
	}

	for _, city := range json {

		obj := dbal.City{Id: city["id"].(int), Name: city["name"].(string), Regexp: city["regexp"].(string), HasSubway: city["has_subway"].(bool)}

		db.AddCity(obj)
	}
}

func loadSubways(db *dbal.DBAL, path string) {
	data, read_err := ioutil.ReadFile(path)

	if read_err != nil {
		fmt.Println(read_err)
		os.Exit(1)
	}

	json := make([]map[string]interface{}, 0)

	parse_err := yaml.Unmarshal(data, &json)
	if parse_err != nil {
		fmt.Println(parse_err)
		os.Exit(1)
	}

	for _, subway := range json {

		obj := dbal.Subway{Id: subway["id"].(int), Name: subway["name"].(string), Regexp: subway["regexp"].(string), City: subway["city"].(int)}

		db.AddSubway(obj)
	}
}

func main() {

	args := os.Args
	if len(args) < 3 {
		fmt.Println("Fixture path arg is required")
		os.Exit(1)
	}

	path := args[2]

	db := connectDb()

	fmt.Println("Loading fitures...")

	loadCities(db, fmt.Sprintf("%s/city.yml", path))
	loadSubways(db, fmt.Sprintf("%s/subway_ekaterinburg.yml", path))
	loadSubways(db, fmt.Sprintf("%s/subway_kazan.yml", path))
	loadSubways(db, fmt.Sprintf("%s/subway_moskva.yml", path))
	loadSubways(db, fmt.Sprintf("%s/subway_nizny_novgorod.yml", path))
	loadSubways(db, fmt.Sprintf("%s/subway_samara.yml", path))
	loadSubways(db, fmt.Sprintf("%s/subway_sankt_peterburg.yml", path))

	fmt.Println("Loading fitures... done")
}
