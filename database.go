package main

import (
	"log"
	"sync"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Network contains data about one network.
type Network struct {
	BSSID         string  `json:"bssid"`
	SSID          string  `json:"ssid"`
	Frequency     int     `json:"frequency"`
	Capabilities  string  `json:"capabilities"`
	LastTime      int     `json:"lasttime"`
	LastLatitude  float64 `json:"lastlat"`
	LastLongitude float64 `json:"lastlon"`
	Type          string  `json:"type"`
	BestLevel     int     `json:"bestlevel"`
	BestLatitude  float64 `json:"bestlat"`
	BestLongitude float64 `json:"bestlon"`
}

type database struct {
	sync.RWMutex
	path string
	conn *sql.DB
}

func newDatabase(path string) (*database, error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return &database{
		path: path,
		conn: conn,
	}, nil
}

const (
	queryString = "SELECT bssid, ssid, frequency, capabilities, lasttime, lastlat, lastlon, type, bestlevel, bestlat, bestlon FROM network WHERE type = 'W' LIMIT ? OFFSET ?"
	totalQuery  = "SELECT COUNT(bssid) FROM network WHERE type = 'W'"
)

func (db *database) Count() (int, error) {
	db.Lock()
	defer db.Unlock()

	rows, err := db.conn.Query(totalQuery)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, nil
	}

	var count int
	if err := rows.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (db *database) Query(offset, limit int) ([]Network, error) {
	db.Lock()
	defer db.Unlock()

	rows, err := db.conn.Query(queryString, limit, offset)
	if err != nil {
		return []Network{}, err
	}
	defer rows.Close()

	var networks []Network
	for rows.Next() {
		net := Network{}
		err := rows.Scan(
			&net.BSSID,
			&net.SSID,
			&net.Frequency,
			&net.Capabilities,
			&net.LastTime,
			&net.LastLatitude,
			&net.LastLongitude,
			&net.Type,
			&net.BestLevel,
			&net.BestLatitude,
			&net.BestLongitude)
		if err != nil {
			log.Printf("Error: %s", err)
			break
		}

		networks = append(networks, net)
	}

	return networks, nil
}
