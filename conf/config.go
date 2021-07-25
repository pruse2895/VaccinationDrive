package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

type Config struct {

	// APPLICATION
	APP_NAME string `json:"app_name"`

	// SERVER CONFIG
	ADDRESS string `json:"http_address"` // http address to listern  Eg : http://localhost
	PORT    int    `json:"port"`

	// DATABASE CONFIG
	DB_TYPE                  string `json:"type"`
	DB_NAME                  string `json:"db_name"`
	DB_USERNAME              string `json:"username"`
	DB_PASSWORD              string `json:"password"`
	DB_ADDRESS               string `json:"db_address"`
	Max_Connection_Pool_Size int    `json:"max_connection_pool_size"`
}

var (
	Cfg  Config
	once sync.Once
)

// Parse parses the json configuration file
// And converting it into native type
func Parse(file string) error {
	once.Do(func() {
		// Reading the flags
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Println("Error in ReadFile:", err)
		}
		if err := json.Unmarshal(data, &Cfg); err != nil {
			log.Println("Error in Unmarshal:", err)
		}
	})
	return nil
}
