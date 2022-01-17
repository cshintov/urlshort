package main

import (
	"fmt"
    "flag"
    "io/ioutil"
	"net/http"

	"urlshortner/urlshort"
)

func main() {
    var urlsFile, format string
    var handler http.Handler

    flag.StringVar(&format, "fmt", "yaml", "fmt of the urls file")
    flag.StringVar(&urlsFile, "urls", "", "urls as a yaml/json file")

    flag.Parse()

    urls, err := ioutil.ReadFile(urlsFile)
    if err != nil {
        panic(err)
    }

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

    var handlerFunc (func([]byte, http.Handler) (http.HandlerFunc, error))

    switch format {
    case "yaml":
        handlerFunc = urlshort.YAMLHandler

    case "json":
        handlerFunc = urlshort.JSONHandler

    default:
        handlerFunc = nil
        handler = mapHandler
    }

    if handlerFunc != nil {
        handler, err = handlerFunc([]byte(urls), mapHandler)
        if err != nil {
            panic(err)
        }
    }

    fmt.Println("Starting the server on :3000")
    http.ListenAndServe(":3000", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
