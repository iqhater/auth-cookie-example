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

	r := &Routes{
		routes: []string{"/", "/login", "/logout", "/error", "/user"},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", r.notFound(login))
	mux.HandleFunc("/login", u.validate(s.isAuth(admin)))
	mux.HandleFunc("/logout", s.isAuth(s.logout))
	mux.HandleFunc("/user", s.isAuth(admin))
	mux.HandleFunc("/error", unAuth)
	mux.HandleFunc("/404", pageNotFound)

	mux.Handle("./favicon.ico", http.NotFoundHandler())

	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// Redirect HTTP requests to HTTPS
	go http.ListenAndServe(":8080", showLog(http.HandlerFunc(redirectToHTTPS)))

	log.Fatal(http.ListenAndServeTLS(":4433", "tls/auth.signin.dev+1.pem", "tls/auth.signin.dev+1-key.pem", showLog(secureHeaders(mux))))
}
