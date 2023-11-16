package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	sm := http.NewServeMux()
	sm.HandleFunc("/", helloWorldHandler)
	l, err := net.Listen("tcp4", ":8081")
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.Serve(l, sm))
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("------ REQUEST START -----")
	printHeadersToConsole(r)
	saveToFile(r)
	fmt.Println("------ REQUEST END -----")
	w.WriteHeader(200)
	io.WriteString(w, "Hello world!\r\n")
}

func printHeadersToConsole(r *http.Request) {
	fmt.Println("- " + r.Method + " " + r.RequestURI + " (" + r.Proto + ")")
	fmt.Println("- " + r.Host + " (" + r.RemoteAddr + ")")
	fmt.Println("-- HEADERS START --")
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Println("-", name+":", value)
		}
	}
	fmt.Println("-- HEADERS END --")
}

func saveToFile(r *http.Request) {
	file := createFile()
	file.WriteString("+++++++++ REQUEST +++++++++\n")
	file.WriteString("- " + r.Method + " " + r.RequestURI + " (" + r.Proto + ")\n")
	file.WriteString("- " + r.Host + " (" + r.RemoteAddr + ")\n")
	file.WriteString("+++++++++ HEADERS +++++++++\n")
	for name, values := range r.Header {
		// Loop over all values for the name.
		for _, value := range values {
			_, err3 := file.WriteString(name + ": " + value + "\n")
			if err3 != nil {
				log.Fatal(err3)
			}
		}
	}
	file.WriteString("+++++++++ PAYLOAD +++++++++\n")
	b, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	_, err2 := file.WriteString(string(b))
	if err2 != nil {
		log.Fatal(err2)
	}
	defer file.Close()
	fmt.Println("-- SAVING BODY END --")
}

func createFile() *os.File {
	sec := time.Now().Unix()
	fileName := "/tmp/request-" + strconv.FormatInt(sec, 10) + ".txt"
	fmt.Println("- " + fileName)
	f, err := os.Create(fileName)

	if err != nil {
		log.Fatal(err)
	}
	return f
}
