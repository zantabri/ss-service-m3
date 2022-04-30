package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/zantabri/ss-service/handlers"
	"github.com/zantabri/ss-service/store"
)

func main() {

	var sd = flag.String("sd", "", "secret directory required")
	flag.Parse()

	if len(*sd) == 0 {
		panic("sd : secret directory is required")
	}

	file_store, err := store.NewFileStore(sd)

	if err != nil {
		panic(err.Error())
	}

	salt := os.Getenv("SS_SALT")
	password := os.Getenv("SS_PASSWORD")

	if len(salt) == 0 || len(password) == 0 {
		panic("either salt or password not set")
	}

	enc_store, err := store.NewEncryptedFileStore(file_store, salt, password)

	if err != nil {
		panic(err.Error())
	}

	handlers := handlers.New(&enc_store)
	router := httprouter.New()

	if err != nil {
		panic(err.Error())
	}

	router.GET("/health", handlers.HealthCheck)
	router.POST("/", handlers.AddSecret)
	router.GET("/", handlers.GetSecret)
	http.ListenAndServe(":8080", router)

}
