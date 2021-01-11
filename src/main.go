package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"net"
	"flag"
	"unit.nginx.org/go"
)

func page(colour string, ip_address string, user_agent string, vorteil_cloud string, server_ip string) string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
 <meta charset="UTF-8">
 <title>Hello World</title>
 <meta name="description" content="Description">
 <meta name="viewport" content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
 <meta http-equiv="refresh" content="2">
 <link rel="stylesheet" href="https://www.w3schools.com/w3css/4/w3.css">
 <!-- Compiled and minified CSS -->
 <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/materialize/0.100.2/css/materialize.min.css">
 <!-- Compiled and minified JavaScript -->
 <script src="https://cdnjs.cloudflare.com/ajax/libs/materialize/0.100.2/js/materialize.min.js"></script>
 <style>
 body {
	 background-color: ` + colour + `;
 }
 h5 {
	 color: #2396d8;
 }
 </style>
</head>
<body  class="valign-wrapper" style="height:100vh;">
<div class="row">
<div class="center-align">
<img src="data:image/png;base64,` + picture + `" alt="Vorteil">
<h5>WELCOME TO VORTEIL</h5>
</div>
</div>
<div class="w3-container w3-content w3-center w3-padding-64" style="max-width:800px">
<h2 class="w3-wide">CONNECTION INFORMATION</h2>
<p class="w3-opacity"><i>Visitor Information</i></p>
<p class="w3-justify">Visitor IP address: ` + ip_address + `</p>
<p class="w3-justify">Visitor User agent: ` + user_agent + `</p>
<p></p>
<p class="w3-justify">Server IP address: ` + server_ip + `</p>
<p class="w3-justify">Server hosted on: ` + vorteil_cloud + `</p>
<p class="w3-justify">Server local time: ` + time.Now().String() + `</p>
<p></p>
</div>
</body>
</html>`
}

func getRealAddr(r *http.Request)  string {

    remoteIP := ""
    // the default is the originating ip. but we try to find better options because this is almost
    // never the right IP
    if parts := strings.Split(r.RemoteAddr, ":"); len(parts) == 2 {
        remoteIP = parts[0]
    }
    // If we have a forwarded-for header, take the address from there
    if xff := strings.Trim(r.Header.Get("X-Forwarded-For"), ","); len(xff) > 0 {
        addrs := strings.Split(xff, ",")
        lastFwd := addrs[len(addrs)-1]
        if ip := net.ParseIP(lastFwd); ip != nil {
            remoteIP = ip.String()
        }
    // parse X-Real-Ip header
    } else if xri := r.Header.Get("X-Real-Ip"); len(xri) > 0 {
        if ip := net.ParseIP(xri); ip != nil {
            remoteIP = ip.String()
        }
    }

    return remoteIP

}

const white = "#FFFFFF"

func main() {
	//colour := os.Getenv("BACKGROUND")
	cloud := os.Getenv("CLOUD_PROVIDER")
	//colour := os.Args[1]
	wordPtr := flag.String("colour","FFFFFF","a string")
	flag.Parse()
	colour := *wordPtr

	if colour == "" {
		log.Printf("No background color set in BACKGROUND environment variable\n")
		colour = white
	}

	colour = "#" + strings.TrimPrefix(strings.TrimPrefix(colour, "0x"), "#")
	colour = strings.ToUpper(colour)

	if len(colour) != 7 {
		log.Printf("Invalid BACKGROUND color: must be six characters of hexadecimal (like '0xFFFFFF')\n")
		colour = white
	}

	valid := true
	for i := 1; i < len(colour); i++ {
		c := colour[i]
		if c < '0' || c > 'F' || (c > '9' && c < 'A') {
			valid = false
		}
	}
	if !valid {
		log.Printf("Invalid BACKGROUND color: non-hexadecimal characters detected\n")
		colour = white
	}

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("SERVING REQUEST")
		addr := getRealAddr(r)
		agent := r.Header.Get("User-Agent")
		serverip := "UKNOWN"
		a, ok := r.Context().Value(http.LocalAddrContextKey).(net.Addr)
			if !ok {
			// handle address not found
		} else {
			ta, ok := a.(*net.TCPAddr)
			if !ok {
				// handle unknown address type
			} else {
				serverip = ta.IP.String()
			}
		}
		rdr := strings.NewReader(page(colour, addr, agent, cloud, serverip))
		_, err := io.Copy(w, rdr)
		if err != nil {
			log.Printf("Connection error: %v\n", err.Error())
		}
	})
	port := os.Getenv("BIND")
	if port == "" {
		port = "80"
	}
	log.Printf("Binding port: %s\n", port)
	err = unit.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
