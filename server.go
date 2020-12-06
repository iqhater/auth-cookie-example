package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

// init is invoked before main()
func init() {

	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {

	// init config environments
	cfg := NewConfig()

	// init user structure
	u := NewUser(cfg)

	// init session structure
	s := NewSession()

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
	go http.ListenAndServe(":"+cfg.httpPort, showLog(http.HandlerFunc(redirectToHTTPS)))

	log.Fatal(http.ListenAndServeTLS(":"+cfg.httpsPort, cfg.tlsCert, cfg.tlsKey, showLog(secureHeaders(mux))))
}
