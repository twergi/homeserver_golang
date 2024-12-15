package main

import (
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"
)

var IP string = getLocalIP()
var Port string = "80"

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr).IP
	return localAddr.String()
}

func diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("static/html/index.html")

	loc, _ := time.LoadLocation("Europe/Moscow")
	start := time.Date(2018, 1, 9, 18, 0, 0, 0, loc)
	year, month, day, hour, min, sec := diff(start, time.Now())

	tmpl.Execute(w, struct {
		HostIP string
		Year   int
		Month  int
		Day    int
		Hour   int
		Min    int
		Sec    int
	}{
		HostIP: IP,
		Year:   year,
		Month:  month,
		Day:    day,
		Hour:   hour,
		Min:    min,
		Sec:    sec,
	})
}

func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("static/html/shutdown.html")

	if err != nil {
		panic(err)
	}
	defer file.Close()

	data := make([]byte, 512)
	_, err = file.Read(data)

	if err != nil {
		panic(err)
	}

	w.Write(data)

	cmd := exec.Command("shutdown", "now")

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func restartHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("static/html/restart.html")

	if err != nil {
		panic(err)
	}
	defer file.Close()

	data := make([]byte, 512)
	_, err = file.Read(data)

	if err != nil {
		panic(err)
	}

	w.Write(data)

	cmd := exec.Command("shutdown", "-r", "now")

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/shutdown", shutdownHandler)
	http.HandleFunc("/restart", restartHandler)
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer((http.Dir("static/css")))))

	fullAddr := fmt.Sprintf("%s:%s", IP, Port)

	println("starting server on:", fullAddr)
	err := http.ListenAndServe(fullAddr, nil)

	if err != nil {
		panic(err)
	}
}
