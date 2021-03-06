package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type QueryResult struct {
	Total    int             `json:"total"`
	Networks []NetworkMarker `json:"networks"`
}

// NetworkMarker contains the information for a map marker.
type NetworkMarker struct {
	Title         string  `json:"title"`
	Latitude      float64 `json:"lat"`
	Longitude     float64 `json:"lon"`
	Text          string  `json:"text"`
	Secure        bool    `json:"secure"`
	HighFrequency bool    `json:"highFreq"`
}

var markerInfoTmpl = `<table class="table network">
  <tr>
    <th>SSID</th>
    <td>{{.SSID}}</td>
  </tr>
  <tr>
    <th>BSSID</th>
    <td>{{.BSSID}}</td>
  </tr>
  <tr>
    <th>Frequency</th>
    <td>{{.Frequency}}</td>
  </tr>
  <tr>
    <th>Capabilities</th>
    <td>
      {{$caps := .Capabilities | splitCapabilities}}
      {{range $caps}}
      <span class="badge">{{.}}</span>
      {{end}}
    </td>
  </tr>
  <tr>
    <th>Last seen</th>
    <td>{{.LastTime | time}}</td>
  </tr>
  <tr>
    <th>Type</th>
    <td>{{.Type}}</td>
  </tr>
</table>`

const defaultLimit = "500"

func formatTime(timeEpochMillis int) string {
	t := time.Unix(int64(timeEpochMillis/1000), 0)
	return t.UTC().Format(time.RFC3339)
}

func splitCapabilities(capabilities string) []string {
	tokens := strings.Split(capabilities, "][")
	caps := make([]string, 0, len(tokens))
	for _, t := range tokens {
		caps = append(caps, strings.Trim(t, "[]"))
	}
	return caps
}

func QueryHandler(db *database) http.Handler {
	templateFuncs := map[string]interface{}{
		"time":              formatTime,
		"splitCapabilities": splitCapabilities,
	}
	markerInfo, err := template.New("marker").Funcs(templateFuncs).Parse(markerInfoTmpl)
	if err != nil {
		log.Fatalf("Error parsing template: %s", err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limitStr := r.URL.Query().Get("limit")
		if len(limitStr) == 0 {
			limitStr = defaultLimit
		}

		offsetStr := r.URL.Query().Get("offset")
		if len(offsetStr) == 0 {
			offsetStr = "0"
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		networks, err := db.Query(offset, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		markers := make([]NetworkMarker, 0, len(networks))
		for _, net := range networks {
			text := &bytes.Buffer{}
			if err := markerInfo.Execute(text, net); err != nil {
				log.Printf("Error executing template: %s", err)
				continue
			}

			secure := strings.Contains(net.Capabilities, "[WPA2")

			marker := NetworkMarker{
				Title:         net.SSID,
				Latitude:      net.BestLatitude,
				Longitude:     net.BestLongitude,
				Text:          text.String(),
				Secure:        secure,
				HighFrequency: net.Frequency > 5000,
			}
			markers = append(markers, marker)
		}

		total, err := db.Count()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result := QueryResult{
			Total:    total,
			Networks: markers,
		}

		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
