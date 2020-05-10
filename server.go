package main

import (
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

func main() {

	user, exists := os.LookupEnv("LOGIN")
	if !exists {
		log.Println("Env variable LOGIN does not exist!")
	}

	pass, exists := os.LookupEnv("PASSWORD")
	if !exists {
		log.Println("Env variable PASSWORD does not exist!")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}

	u := &User{
		login:        user,
		password:     pass,
		passwordHash: passwordHash,
	}

	s := &Session{
		ID:   "",
		Name: "",
	}

	http.HandleFunc("/", notFound(login))
	http.HandleFunc("/login", showLog(u.validate(s.isAuth(admin))))
	http.HandleFunc("/logout", showLog(s.isAuth(s.logout)))
	http.HandleFunc("/user", showLog(s.isAuth(admin)))
	http.HandleFunc("/error", showLog(unAuth))
	http.HandleFunc("/404", showLog(pageNotFound))

	http.Handle("./favicon.ico", http.NotFoundHandler())

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	log.Fatal(http.ListenAndServe(":4450", nil))
}
