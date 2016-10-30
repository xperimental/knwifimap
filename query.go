package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

// NetworkMarker contains the information for a map marker.
type NetworkMarker struct {
	Title     string  `json:"title"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Text      string  `json:"text"`
}

var markerInfoTmpl = `<div class="network">
  <div class="ssid">SSID: {{.SSID}}</div>
  <div class="bssid">BSSID: {{.BSSID}}</div>
  <div class="frequency">Frequency: {{.Frequency}}</div>
  <div class="capabilities">Capabilities: {{.Capabilities}}</div>
  <div class="lastseen">Last seen: {{.LastTime}}</div>
  <div class="type">Type: {{.Type}}</div>
</div>`

func QueryHandler(db *database) http.Handler {
	markerInfo, err := template.New("marker").Parse(markerInfoTmpl)
	if err != nil {
		log.Fatalf("Error parsing template: %s", err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		networks, err := db.Query()
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

			marker := NetworkMarker{
				Title:     net.SSID,
				Latitude:  net.BestLatitude,
				Longitude: net.BestLongitude,
				Text:      text.String(),
			}
			markers = append(markers, marker)
		}

		if err := json.NewEncoder(w).Encode(markers); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
