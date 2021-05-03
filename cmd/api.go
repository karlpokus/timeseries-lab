package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	"timeseries/lib/db"
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

func search(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sr db.SearchRequest
		if err := json.NewDecoder(r.Body).Decode(&sr); err != nil {
			fail(w, err)
			return
		}
		defer r.Body.Close()
		log.Printf("POST /search for type: %s target: %s", sr.Type, sr.Target)
		keys, err := db.Keys(pool)
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

func query(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var q db.QueryRequest
		if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
			fail(w, err)
			return
		}
		defer r.Body.Close()
		log.Printf("POST /query for %s and %s", q.From, q.To)
		log.Printf("AdhocFilters %d", len(q.AdhocFilters))
		log.Printf("Targets %+v", q.Targets)
		log.Printf("MaxDataPoints %d", q.MaxDataPoints)
		rcds, err := db.Query(pool, q)
		if err != nil {
			fail(w, err)
			return
		}
		log.Printf("found %d records", len(rcds))
		type Response struct {
			Target     string          `json:"target"`
			Datapoints [][]interface{} `json:"datapoints"`
		}
		// quick and dirty
		var out []Response
		bat := Response{
			Target:     "battery",
			Datapoints: telemetry.Datapoints(rcds, "battery"),
		}
		heat := Response{
			Target:     "heat",
			Datapoints: telemetry.Datapoints(rcds, "heat"),
		}
		hog := Response{
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
	url := "postgres://postgres:secret@pg:5432/test"
	pool, err := db.Connect(url)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("db connected")
	defer pool.Close()

	http.HandleFunc("/api/v1", ok)                      // GET
	http.HandleFunc("/api/v1/search", search(pool))     // POST
	http.HandleFunc("/api/v1/query", query(pool))       // POST
	http.HandleFunc("/api/v1/annotations", annotations) // POST
	log.Println("v1 running on port 8989")
	log.Fatal(http.ListenAndServe(":8989", nil))
}
