package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"github.com/stianeikeland/go-rpio"
)

var (
	quit chan bool
)

var log = logrus.New()

const ()

func init() {
	log.Formatter = new(logrus.JSONFormatter)
	log.Formatter = new(logrus.TextFormatter) // default
	log.Level = logrus.DebugLevel
}

func main() {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	defer func() {
		err := recover()
		if err != nil {
			log.WithFields(logrus.Fields{
				"omg":    true,
				"err":    err,
				"number": 100,
			}).Fatal("The ice breaks!")
		}
	}()

	// command line flags
	viper.SetDefault("port", 9999)

	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// setup routes
	router := mux.NewRouter()

	router.Path("/api/status").Methods("GET").HandlerFunc(getGPIOStatus).Methods("GET")
	router.Path("/api/pins").Methods("GET").HandlerFunc(getPins).Methods("GET")
	router.Path("/api/gpio/{id}").Methods("GET").HandlerFunc(getStatus).Methods("GET")
	router.Path("/api/gpio/{id}/toggle").Methods("POST").HandlerFunc(toggle)
	router.Path("/api/gpio/{id}/high").Methods("POST").HandlerFunc(high)
	router.Path("/api/gpio/{id}/low").Methods("POST").HandlerFunc(low)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	http.Handle("/api/", handler)
	http.Handle("/", http.StripPrefix("/", http.FileServer(assetFS())))

	log.Printf("Running on port %d\n", viper.GetInt("port"))

	addr := fmt.Sprintf("0.0.0.0:%d", viper.GetInt("port"))

	go func() {
		err := http.ListenAndServe(addr, nil)
		fmt.Println(err.Error())
	}()

	select {
	case <-quit:
		os.Exit(0)
	}
}
