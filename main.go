package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"time"
	"vaccinationDrive/conf"
	"vaccinationDrive/dbcon"
	"vaccinationDrive/dbscripts"
	"vaccinationDrive/routes"

	"github.com/rs/cors"
)

func main() {
	log.Println(time.Now().UTC())

	var configFile = flag.String("conf", "", "configuration file(mandatory)")

	flag.Parse()
	if flag.NFlag() != 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Parsing configuration
	if err := conf.Parse(*configFile); err != nil {
		log.Fatalln("ERROR: ", err)
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cpu, _ := strconv.Atoi(os.Getenv("GOMAXPROCS"))
	runtime.GOMAXPROCS(cpu)
	log.Println("INFO: Number of cpu configured - ", cpu)

	dbcon.Connect()
	defer dbcon.Close()

	dbscripts.InitDB()

	router := routes.RouterConfig()
	//r := chi.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders: []string{"Origin", "X-Requested-With", "Content-Type", "Accept",
			"Authorization", "Access-Control-Allow-Headers", "Access-Control-Allow-Origin"},
	})

	server := http.Server{
		Addr:         fmt.Sprintf(":%d", conf.Cfg.PORT),
		ReadTimeout:  90 * time.Second,
		WriteTimeout: 90 * time.Second,
		Handler:      c.Handler(router),
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	//Graceful shut down
	go func() {
		<-quit
		log.Println("Server is shutting down...")

		//Close resources before shut down
		dbcon.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		//Shutdown server
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Unable to gracefully shutdown the server: %v\n", err)
		}

		//Close channels
		close(quit)
		close(done)
	}()

	log.Printf("Listening on: %d", conf.Cfg.PORT)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Error in listening server: %s", err.Error())
	}
	<-done
	log.Fatal("Server stopped")
}
