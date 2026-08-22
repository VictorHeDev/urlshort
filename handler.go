package urlshort

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...

	return func(w http.ResponseWriter, r *http.Request) {
		// map any key to their corresponding URL
		if path, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, path, http.StatusFound)
			return
		}

		// else return fallback http.Handler
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// TODO: Implement this...
	type pathURL struct {
		Path string `yaml:"path"`
		URL  string `yaml:"url"`
	}

	var pathURLs []pathURL
	err := yaml.Unmarshal(yml, &pathURLs)
	if err != nil {
		return nil, err
	}

	urlMap := make(map[string]string, len(pathURLs))
	for _, v := range pathURLs {
		urlMap[v.Path] = v.URL
	}

	return MapHandler(urlMap, fallback), nil
}

// JSONHandler will parse the provided JSON and return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
//	{
//	  "/some-path": "https://www.some-url.com/demo"
//	}
//
// The only errors that can be returned all related to having
// invalid JSON data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	urlMap := make(map[string]string)
	err := json.Unmarshal(jsn, &urlMap)
	if err != nil {
		return nil, err
	}

	return MapHandler(urlMap, fallback), nil
}

// FileHandler takes the file's suffix passed in via CLI flag
// and returns the appropriate handler
func FileHandler(filename string, fallback http.Handler) (http.HandlerFunc, error) {
	fileSuffix := filepath.Ext(filename)
	byteBody, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	switch fileSuffix {
	case ".json":
		return JSONHandler(byteBody, fallback)
	case ".yaml", ".yml":
		return YAMLHandler(byteBody, fallback)
	default:
		return nil, fmt.Errorf("unsupported filetype for filename: %s", filename)
	}
}
