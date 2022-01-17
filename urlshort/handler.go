package urlshort

import (
	"fmt"
	"net/http"

    "github.com/go-yaml/yaml"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if url, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, url, http.StatusMovedPermanently)
		} else {
			fmt.Println(url, "Map not found!")
			fallback.ServeHTTP(w, r)
		}
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYaml(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

type Url struct {
	Key    string `yaml:"path"`
	Target string `yaml:"url"`
}

func parseYaml(yml []byte) ([]Url, error) {
    var urls []Url
    err := yaml.Unmarshal(yml, &urls)
	if err != nil {
		return nil, err
	}
    return urls, nil
}

func buildMap(urls []Url) map[string]string {
    var data = map[string]string{}
    for _, url := range urls {
        data[url.Key] = url.Target
    }
    return data
}
