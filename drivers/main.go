package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Driver é a estrutura básica de um motorista
type Driver struct {
	UUID  string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

//Drivers contém um slice de Driver
type Drivers struct {
	Drivers []Driver
}

// getDrivers retorna um []byte contendo todos os motoristas
func getDrivers() []byte {
	file, err := os.Open("./data.json")
	if err != nil {
		log.Fatalln("failed to open data.json", err)
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln("failed to read data.json", err)
	}

	return data
}

// ShowDrivers escreve na response a lista de motoristas
func ShowDrivers(resp http.ResponseWriter, req *http.Request) {
	drivers := getDrivers()
	resp.Write([]byte(drivers))
}

// ShowDrivers escreve na response a lista de motoristas
func GetDriverByID(resp http.ResponseWriter, req *http.Request) {
	//buscar os parametros via query parm
	param := mux.Vars(req)
	data := getDrivers()

	//decode para struct
	var drivers Drivers
	json.Unmarshal(data, &drivers)

	//consulta
	for _, d := range drivers.Drivers {
		if d.UUID == param["id"] {
			//encoding da struct para json
			driver, _ := json.Marshal(d)
			resp.Write([]byte(driver))
			return
		}
	}

	resp.Write([]byte("not found"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/drivers", ShowDrivers)
	r.HandleFunc("/drivers/{id}", GetDriverByID)

	http.ListenAndServe(":3001", r)
}
