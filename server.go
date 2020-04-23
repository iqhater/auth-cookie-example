package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

var LOGIN string
var PASSWORD string
var passwordHash []byte

func main() {

	user, exists := os.LookupEnv("LOGIN")
	if !exists {
		fmt.Println("Env variable LOGIN does not exist!")
	}

	LOGIN = user

	pass, exists := os.LookupEnv("PASSWORD")
	if !exists {
		fmt.Println("Env variable PASSWORD does not exist!")
	}

	PASSWORD = pass

	passwordHash, _ = bcrypt.GenerateFromPassword([]byte(PASSWORD), bcrypt.MinCost)

	http.HandleFunc("/", notFound(login))
	http.HandleFunc("/login", showLog(validate(isAuth(admin))))
	http.HandleFunc("/logout", showLog(isAuth(logout)))
	http.HandleFunc("/user", showLog(isAuth(admin)))
	http.HandleFunc("/error", showLog(unAuth))
	http.HandleFunc("/404", showLog(pageNotFound))
	http.Handle("./favicon.ico", http.NotFoundHandler())

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	log.Fatal(http.ListenAndServe(":4450", nil))
}
