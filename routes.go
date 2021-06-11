package main

// Routes struct is used for list of endpoints urls
type Routes struct {
	*Session
	routes []string // slice of web app routes
	files  []string // slice of file names. Used for list of secured public files
}
