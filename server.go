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
		log.Println("No .env file found")
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

	httpPort, exists := os.LookupEnv("HTTP_PORT")
	if !exists {
		log.Println("Env variable HTTP_PORT does not exist!")
	}

	httpsPort, exists := os.LookupEnv("HTTPS_PORT")
	if !exists {
		log.Println("Env variable HTTPS_PORT does not exist!")
	}

	tlsCert, exists := os.LookupEnv("TLS_CERT_PATH")
	if !exists {
		log.Println("Env variable TLS_CERT_PATH does not exist!")
	}

	tlsKey, exists := os.LookupEnv("TLS_KEY_PATH")
	if !exists {
		log.Println("Env variable TLS_KEY_PATH does not exist!")
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

	// list of all routes in app
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
	go http.ListenAndServe(":"+httpPort, showLog(http.HandlerFunc(redirectToHTTPS)))

	log.Fatal(http.ListenAndServeTLS(":"+httpsPort, tlsCert, tlsKey, showLog(secureHeaders(mux))))
}
