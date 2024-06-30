package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var PORT int = 8080
var VERBOSE bool = false
var USE_PROXY_HEADER bool = false // Trust the proxy forwarded header
const VERSION string = "24.7.0"

// Parse header and try to get IP-address.
// Does not actually check IP-address validity in any way.
func getIPFromHeader(headerVal []string) (string, error) {
	if len(headerVal) < 1 {
		return "", fmt.Errorf("no IP found in proxy header")
	}

	ip := strings.Join(headerVal, " ")
	splitIndex := strings.LastIndex(ip, ",") // If header contained multiple addresses

	if splitIndex > 1 {
		ip = ip[splitIndex+1:]
	}

	return strings.TrimSpace(ip), nil
}

// Checks if headers contain proxy address keys
// found key returned or "" and error
func isProxyHeaderSet(headers map[string][]string) (string, error) {
	PROXY_HEADER_KEYS := [2]string{"X-Real-Ip", "X-Forwarded-For"}

	for i := 0; i < 2; i++ {
		_, exists := headers[PROXY_HEADER_KEYS[i]]
		if exists {
			return PROXY_HEADER_KEYS[i], nil
		}
	}
	return "", fmt.Errorf("proxy header key not found in headers")
}

func checkipHandler(resW http.ResponseWriter, req *http.Request) {
	var logBuilder strings.Builder
	var clientIP = ""
	var clientPort = ""

	fmt.Fprintf(&logBuilder, "A %q Request to %q %q from %q. ", req.Method, req.Host, req.RequestURI, req.RemoteAddr)

	if USE_PROXY_HEADER {
		headerKey, err := isProxyHeaderSet(req.Header)
		if err != nil {
			log.Printf("WARNING, Use proxy set, but %s", err)
		}

		if headerKey != "" {
			clientIP, err = getIPFromHeader(req.Header[headerKey])
			if err != nil {
				log.Printf("WARNING, Found header %q, but %s", headerKey, err)
			}

			fmt.Fprintf(&logBuilder, "\"%s: %s\" ", headerKey, clientIP)
		}
	}

	if clientIP == "" { // Did not get IP from a header
		splitIndex := strings.LastIndex(req.RemoteAddr, ":") // split address and port
		if splitIndex > 0 {
			clientIP = req.RemoteAddr[:splitIndex]
			clientPort = req.RemoteAddr[splitIndex+1:]
		}
	}

	if VERBOSE { // Log all headers
		for name, headers := range req.Header {
			for _, val := range headers {
				fmt.Fprintf(&logBuilder, "\"%s: %s\" ", name, val)
			}
		}
	}

	log.Print(logBuilder.String())

	fmt.Fprintf(resW, "Your IP is: %s", clientIP)
	if clientPort != "" {
		fmt.Fprintf(resW, " Port: %s", clientPort)
	}
}

func main() {
	flag.IntVar(&PORT, "port", PORT, "Port to use")
	flag.BoolVar(&VERBOSE, "verbose", VERBOSE, "Show more information of request")
	flag.BoolVar(&USE_PROXY_HEADER, "proxy", USE_PROXY_HEADER, "Trust proxy header to get IP-address")
	flag.Parse()

	http.HandleFunc("/", checkipHandler)
	http.HandleFunc("/checkip", checkipHandler)

	log.Printf("GoCheckIP server %v at port %v TCP\n", VERSION, PORT)

	var addrAndPort string = fmt.Sprintf(":%v", PORT)
	if err := http.ListenAndServe(addrAndPort, nil); err != nil {
		log.Fatal("Cannot start server. ", err)
	}
}
