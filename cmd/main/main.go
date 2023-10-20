package main

import (
	"net"
	"net/http"
	"os"
	"os/exec"
)

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
	file, err := os.Open("static/index.html")

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
}

func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("static/shutdown.html")

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
	file, err := os.Open("static/restart.html")

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

	ip := getLocalIP()
	port := "8000"
	fullAddr := ip + ":" + port

	println("starting server on:", fullAddr)
	err := http.ListenAndServe(fullAddr, nil)

	if err != nil {
		panic(err)
	}
}
