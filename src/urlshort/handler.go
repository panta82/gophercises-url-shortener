package urlshort

import (
	"net/http"

	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"urlshort/database"
	. "urlshort/types"
)

func doRedirect(rw http.ResponseWriter, req *http.Request, url string) {
	http.Redirect(rw, req, url, 302)
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		if url, ok := pathsToUrls[req.URL.Path]; ok {
			doRedirect(rw, req, url)
			return
		}

		// Fallback
		fallback.ServeHTTP(rw, req)
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
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var redirects []Redirect

	err := yaml.Unmarshal(yml, &redirects)
	if err != nil {
		return nil, err
	}

	pathsToUrls := make(map[string]string)
	for _, redirect := range redirects {
		pathsToUrls[redirect.Path] = redirect.Url
	}

	return MapHandler(pathsToUrls, fallback), nil
}

// Additional handler using bolt DB
func DatabaseHandler(db database.Database, fallback http.Handler) (http.HandlerFunc, error) {
	return func(rw http.ResponseWriter, req *http.Request) {
		url, err := db.GetUrlForPath(req.URL.Path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())

			rw.WriteHeader(500)
			rw.Write([]byte("Internal server error"))
			return
		}

		if url == "" {
			// No url is found, go to fallback
			fallback.ServeHTTP(rw, req)
			return
		}

		doRedirect(rw, req, url)
	}, nil
}
