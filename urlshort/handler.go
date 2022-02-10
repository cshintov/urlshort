package urlshort

import (
	"fmt"
	"net/http"
	"encoding/json"

    "github.com/go-yaml/yaml"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if url, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, url, http.StatusMovedPermanently)
		} else {
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

func JSONHandler(jsondata []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJson, err := parseJson(jsondata)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedJson)
	return MapHandler(pathMap, fallback), nil
}

type Url struct {
    Key    string `yaml:"path" json:"path"`
    Target string `yaml:"url" json:"url"`
}

func parseJson(yml []byte) ([]Url, error) {
    var urls []Url
	if err := json.Unmarshal(yml, &urls); err != nil {
		return nil, err
	}
    return urls, nil
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
