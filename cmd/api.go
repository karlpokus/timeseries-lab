package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"timeseries/lib/store"
	"timeseries/lib/telemetry"
)

func dump(body io.ReadCloser) error {
	var x map[string]interface{}
	if err := json.NewDecoder(body).Decode(&x); err != nil {
		return err
	}
	defer body.Close()
	log.Printf("%+v", x)
	return nil
}

func fail(w http.ResponseWriter, err error) {
	log.Printf("%s", err)
	http.Error(w, "server error", 500)
}

func ok(w http.ResponseWriter, r *http.Request) {
	log.Printf("got %s request at %s", r.Method, r.URL.Path)
	w.WriteHeader(200)
}

func search(st store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sr telemetry.Search
		if err := json.NewDecoder(r.Body).Decode(&sr); err != nil {
			fail(w, err)
			return
		}
		defer r.Body.Close()
		log.Printf("POST /search for type: %s target: %s", sr.Type, sr.Target)
		keys, err := st.Keys()
		if err != nil {
			fail(w, err)
			return
		}
		b, err := json.Marshal(keys)
		if err != nil {
			fail(w, err)
			return
		}
		w.Write(b)
	}
}

func query(st store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var q telemetry.Query
		if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
			fail(w, err)
			return
		}
		defer r.Body.Close()
		log.Printf("POST /query for %s and %s", q.From, q.To)
		log.Printf("AdhocFilters %d", len(q.AdhocFilters))
		log.Printf("Targets %+v", q.Targets)
		log.Printf("MaxDataPoints %d", q.MaxDataPoints)
		rcds, err := st.Query(q)
		if err != nil {
			fail(w, err)
			return
		}
		log.Printf("found %d records", len(rcds))
		// quick and dirty
		var out []telemetry.Response
		bat := telemetry.Response{
			Target:     "battery",
			Datapoints: telemetry.Datapoints(rcds, "battery"),
		}
		heat := telemetry.Response{
			Target:     "heat",
			Datapoints: telemetry.Datapoints(rcds, "heat"),
		}
		hog := telemetry.Response{
			Target:     "hog",
			Datapoints: telemetry.Datapoints(rcds, "hog"),
		}
		out = append(out, bat, heat, hog)
		b, err := json.Marshal(out)
		if err != nil {
			fail(w, err)
			return
		}
		w.Write(b)
	}
}

func annotations(w http.ResponseWriter, r *http.Request) {
	err := dump(r.Body)
	if err != nil {
		http.Error(w, "heplp", 500)
		return
	}
	response := `
  [
    {
      "text": "text shown in body",
      "title": "Annotation Title",
      "isRegion": true,
      "time": "timestamp",
      "timeEnd": "timestamp",
      "tags": ["tag1"]
    }
  ]
  `
	w.Write([]byte(response))
}

func main() {
	st, err := store.New("postgres://postgres:secret@pg:5432/test")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("db connected")
	defer st.Close()
	http.HandleFunc("/api/v1", ok)
	http.HandleFunc("/api/v1/search", search(st))
	http.HandleFunc("/api/v1/query", query(st))
	http.HandleFunc("/api/v1/annotations", annotations)
	log.Println("v1 running on port 8989")
	log.Fatal(http.ListenAndServe(":8989", nil))
}
