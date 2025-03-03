package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website! \n")
}

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello HTTP!\n")
}
func main() {
	var portFlag = flag.String("port", ":8080", "write port for the server")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/hello", getHello)

	err := http.ListenAndServe(*portFlag, mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
