package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	var portFlag = flag.String("port", "8080", "write port for the server")
	var dirFlag = flag.String("dir", "data", "directory for data storing")
	flag.Parse()

	if _, err := os.Stat(*dirFlag); os.IsNotExist(err) {
		fmt.Println("Given directory does not exist, creating directory...")
		err := os.MkdirAll(*dirFlag, os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to create directory: %v\n", err)
			os.Exit(1)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/hello", getHello)

	server := &http.Server{
		Addr:    ":" + *portFlag,
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("Server is running on port %s\n", *portFlag)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// Ожидание сигнала завершения
	<-stop
	fmt.Println("\nShutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Error during shutdown: %s\n", err)
	} else {
		fmt.Println("Server stopped gracefully")
	}
}
