package main

import (
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"os/exec"
)

var IP string = getLocalIP()
var Port string = "8000"

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr).IP
	return localAddr.String()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("static/html/index.html")

	tmpl.Execute(w, struct {
		HostIP string
	}{
		HostIP: IP,
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
