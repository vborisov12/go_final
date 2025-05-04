package api

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// NextDateHendler обработчик запроса /api/nextdate
// опирается на параметры now, date, repeat
// для обработки исопльзует функцию NewDate из scheduler.go
func NextDateHendler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse("20060102", nowStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid 'now' date format"), http.StatusBadRequest)
			log.Printf("Invalid 'now' date format: %v", err)
			return
		}
	}

	if dateStr == "" {
		http.Error(w, fmt.Sprintf("Missing 'date' parameter"), http.StatusBadRequest)
		log.Print("Missing 'date' parameter")
		return
	}

	if repeat == "" {
		http.Error(w, fmt.Sprintf("Missing 'repeat' parameter"), http.StatusBadRequest)
		log.Print("Missing 'repeat' parameter")
		return
	}

	nexDate, err := NewDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, nexDate)
}
