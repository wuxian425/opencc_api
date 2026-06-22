package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/longbridge/opencc"
)

var s2t, t2s, s2tw, tw2s, s2hk, hk2s, s2twp, tw2sp, t2tw, hk2t, t2hk, t2jp, jp2t, tw2t *opencc.OpenCC

func main() {
	// command-line arguments
	portPtr := flag.Int("port", 1145, "define the port service running on")
	loadPtr := flag.String("dicts", "s2tw", "define dicts opencc load (use '+' to separate)")

	flag.Parse()

	for v := range strings.SplitSeq(*loadPtr, "+") {
		var err error

		switch v {
		case "s2t":
			s2t, err = opencc.New(v)
		case "t2s":
			t2s, err = opencc.New(v)
		case "s2tw":
			s2tw, err = opencc.New(v)
		case "tw2s":
			tw2s, err = opencc.New(v)
		case "s2hk":
			s2hk, err = opencc.New(v)
		case "hk2s":
			hk2s, err = opencc.New(v)
		case "s2twp":
			s2twp, err = opencc.New(v)
		case "tw2sp":
			tw2sp, err = opencc.New(v)
		case "t2tw":
			t2tw, err = opencc.New(v)
		case "hk2t":
			hk2t, err = opencc.New(v)
		case "t2hk":
			t2hk, err = opencc.New(v)
		case "t2jp":
			t2jp, err = opencc.New(v)
		case "jp2t":
			jp2t, err = opencc.New(v)
		case "tw2t":
			tw2t, err = opencc.New(v)
		}

		if err != nil {
			log.Fatal("failed to load dict: ", v)
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
			http.Error(w, `{"error": "missing 'text' parameter"}`, http.StatusBadRequest)
			return
		}

		if dict == "" {
			dict = "s2tw"
		}

		var convertedText string
		var err error

		switch dict {
		case "s2t":
			convertedText, err = s2t.Convert(text)
		case "t2s":
			convertedText, err = t2s.Convert(text)
		case "s2tw":
			convertedText, err = s2tw.Convert(text)
		case "tw2s":
			convertedText, err = tw2s.Convert(text)
		case "s2hk":
			convertedText, err = s2hk.Convert(text)
		case "hk2s":
			convertedText, err = hk2s.Convert(text)
		case "s2twp":
			convertedText, err = s2twp.Convert(text)
		case "tw2sp":
			convertedText, err = tw2sp.Convert(text)
		case "t2tw":
			convertedText, err = t2tw.Convert(text)
		case "hk2t":
			convertedText, err = hk2t.Convert(text)
		case "t2hk":
			convertedText, err = t2hk.Convert(text)
		case "t2jp":
			convertedText, err = t2jp.Convert(text)
		case "jp2t":
			convertedText, err = jp2t.Convert(text)
		case "tw2t":
			convertedText, err = tw2t.Convert(text)
		default:
			http.Error(w, fmt.Sprintf(`{"error": "dict %s doesn't imported."}`, dict), http.StatusBadRequest)
			return
		}

		if err != nil {
			log.Fatal("failed to convert", text)
		}

		w.Write([]byte(convertedText))
	})

	// http listen
	port := fmt.Sprintf(":%d", *portPtr)

	fmt.Println("opencc_api service is running on", port)
	err := http.ListenAndServe(port, r)

	if err != nil {
		log.Fatal("failed to run service")
	}
}
