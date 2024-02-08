/*
Serve is a very simple static file server in go
Usage:

	-p="8100": port to serve on
	-d=".":    the directory of static files to host

Navigating to http://localhost:8100 will display the index.html or directory
listing file.
*/
package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	_, ok1 := os.LookupEnv("BASIC_AUTH_NAME")
	_, ok2 := os.LookupEnv("BASIC_AUTH_PASS")
	if ok1 && ok2 {
		/* great */
	} else {
		panic("provide environment vars BASIC_AUTH_NAME and BASIC_AUTH_PASS")
	}

	directory := "./keep"
	port := flag.String("p", "443", "port to serve on")
	flag.Parse()

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		name, pass, ok := request.BasicAuth()
		if !ok {
			writer.WriteHeader(403)
			writer.Write([]byte("forbidden\n"))
			return
		}

		ok = false
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			validName, _ := os.LookupEnv("BASIC_AUTH_NAME")
			validPass, _ := os.LookupEnv("BASIC_AUTH_PASS")
			ok = validName == name && validPass == pass
			wg.Done()
		}()
		d := time.Duration(10+(int(rand.Float64()*10))) * time.Millisecond
		<-time.After(d)
		wg.Wait()

		if ok {
			http.FileServer(http.Dir(directory)).ServeHTTP(writer, request)
		} else {
			writer.WriteHeader(401)
			writer.Write([]byte("unauthorized\n"))
			return
		}

	})

	log.Printf("Serving %s on HTTP port: %s\n", directory, *port)
	err := http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
