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

	// routes: list of all routes in app
	// files: list of all files in public folder if you want to protect from anonymous users
	r := &Routes{
		Session: s,
		routes:  []string{"/", "/login", "/logout", "/error", "/user"},
		files:   []string{"user.html", "user.css", "gopher_wizard.png", "cookie_logo.svg", "jwt_logo.svg"},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", r.notFound(login))
	mux.HandleFunc("/login", s.setAuthType(u.validate(s.setCookie(s.isAuth(admin)))))
	mux.HandleFunc("/logout", s.isAuth(s.logout))
	mux.HandleFunc("/user", s.isAuth(admin))
	mux.HandleFunc("/error", unAuth)
	mux.HandleFunc("/404", pageNotFound)

	mux.Handle("./favicon.ico", http.NotFoundHandler())
	mux.Handle("/public/", r.secureFiles(http.StripPrefix("/public/", http.FileServer(http.Dir("./public")))))

	if cfg.forcedTLS {
		// Redirect HTTP requests to HTTPS
		go http.ListenAndServe(":"+cfg.httpPort, showLog(http.HandlerFunc(redirectToHTTPS)))

		log.Fatal(http.ListenAndServeTLS(":"+cfg.httpsPort, cfg.tlsCert, cfg.tlsKey, showLog(secureHeaders(mux))))
	} else {
		log.Fatal(http.ListenAndServe(":"+cfg.httpPort, showLog(secureHeaders(mux))))
	}
}
