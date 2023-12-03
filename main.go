package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

var rootDir, port string

func flagParse() {
	flag.StringVar(&rootDir, "d", ".", "the directory for sharing")
	flag.StringVar(&port, "p", "8001", "port number")
	flag.Parse()

	if port == "" {
		log.Fatal("prot number can't be empty")
	}

	if rootDir == "" {
		log.Fatal("directory name can't be empty")
	}
}

func localIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "localhost"
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String()
		}
	}

	return "localhost"
}

func main() {
	flagParse()

	// var zipD = filepath.Join(os.TempDir(), "goshra_zip")
	zipD, err := os.MkdirTemp(os.TempDir(), "goshare_zip_")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		os.RemoveAll(zipD)
	}()

	sv := server{rootDir, os.TempDir(), zipD}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/browse/", http.StatusMovedPermanently)
	})

	http.HandleFunc("/browse/", sv.browse)

	http.HandleFunc("/zip", sv.zip)
	fmt.Printf("serving at http://%s:%s\n", localIp(), port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("\nwhile serving err: %v\n", err)
		os.Exit(1)
	}
}
