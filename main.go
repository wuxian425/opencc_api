package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/longbridge/opencc"
)

var importedDicts = make(map[string]*opencc.OpenCC)

func main() {
	// command-line arguments
	portPtr := flag.Int("port", 1145, "define the port service running on")
	dictsPtr := flag.String("dicts", "s2tw", "define dicts opencc load (use '+' to separate)")

	flag.Parse()

	supportedDicts := []string{
		"s2t", "t2s", "s2tw", "tw2s", "s2hk", "hk2s", "s2twp", "tw2sp",
		"t2tw", "hk2t", "t2hk", "t2jp", "jp2t", "tw2t", "s2hk-finance",
	}

	for v := range strings.SplitSeq(*dictsPtr, "+") {
		if slices.Contains(supportedDicts, v) {
			instance, err := opencc.New(v)

			if err != nil {
				log.Fatalf("failed to load dict '%s' : %v", v, err)
			}

			importedDicts[v] = instance
		} else {
			log.Fatalf("opencc doesn't support dict '%s'", v)
		}
	}

	// web service
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/convert", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		text := queryParams.Get("text")
		dict := queryParams.Get("dict")

		if text == "" {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"error": "missing 'text' parameter"}`, http.StatusBadRequest)
			return
		}

		if dict == "" {
			dict = "s2tw"
		}

		val, ok := importedDicts[dict]
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, fmt.Sprintf(`{"error": "dict '%s' isn't imported."}`, dict), http.StatusBadRequest)
			return
		}

		convertedText, err := val.Convert(text)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			log.Printf("convert error: %s", text)
			http.Error(w, fmt.Sprintf(`{"error": "failed to convert '%s' with '%s' dict"`, text, dict), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(convertedText))
	})

	// http listen
	port := fmt.Sprintf(":%d", *portPtr)

	fmt.Printf("opencc_api service is running on %s", port)
	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatalf("failed to run service: %v", err)
	}
}
