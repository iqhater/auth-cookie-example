package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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

var passwordHash, _ = bcrypt.GenerateFromPassword([]byte(PASSWORD), bcrypt.MinCost)

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

	fmt.Println(LOGIN, PASSWORD)

	http.HandleFunc("/", login)
	http.HandleFunc("/user", loggedUser)
	// http.HandleFunc("/error", badLogin)
	http.Handle("./favicon.ico", http.NotFoundHandler())

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	log.Fatal(http.ListenAndServe(":4450", nil))
}

func login(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {

		err := bcrypt.CompareHashAndPassword(passwordHash, []byte(req.FormValue("password")))
		if req.FormValue("login") == LOGIN && err == nil {
			// http.Redirect(w, req, "/user", http.StatusTemporaryRedirect)
			http.ServeFile(w, req, "./public/user.html")
		} else {

			// http.Error(w, "Wrong login or password!", http.StatusUnauthorized)
			// http.Redirect(w, req, "/error", http.StatusUnauthorized)
			http.ServeFile(w, req, "./public/error.html")
		}

	} else {
		http.ServeFile(w, req, "./public/login.html")
	}
}

func loggedUser(w http.ResponseWriter, req *http.Request) {
	// http.Redirect(w, req, "http://localhost:4450/user.html", http.StatusTemporaryRedirect)

	if req.Method == http.MethodPost {

		fmt.Println(req.FormValue("login"))
		fmt.Println(req.FormValue("password"))

		err := req.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}

		for k, v := range req.Form {
			fmt.Println("key:", k)
			fmt.Println("val:", strings.Join(v, ""))
		}

		// body, err := ioutil.ReadAll(req.Body)
		// if err != nil {
		// 	log.Println("Can't get response body!", err)
		// }
		// fmt.Println(string(body))
		http.ServeFile(w, req, "./public/user.html")
	}

	http.ServeFile(w, req, "./public/user.html")
	// w.Header().Set("Content-Type", "text/html")
	// w.Write([]byte("You're not logged in!"))
}

/* func badLogin(w http.ResponseWriter, req *http.Request) {
	// http.Redirect(w, req, "http://localhost:4450/error.html", http.StatusTemporaryRedirect)
	http.ServeFile(w, req, "./public/error.html")
} */
