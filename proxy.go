package main

import (
  "fmt"
  "io"
  "log"
  "net/http"
  "os"
  "time"
)

func startProxyServer() *http.Server {
  http.HandleFunc("/", handleRequest)

  server := &http.Server{
    Addr: ":8080",
  }
  
  return server
}

func handleRequest(writer http.ResponseWriter, req *http.Request) {
	startReq := time.Now()

	targetUrl := os.Getenv("PROXY_TARGET_URL") + fmt.Sprintf("%s", req.URL)
	if targetUrl == "" {
		log.Fatal("PROXY_TARGET_URL in env not setted")
	}

	fmt.Println(req.Method, req.URL)
	
	proxyReq, err := http.NewRequest(req.Method, targetUrl, req.Body)
	if err != nil {
		http.Error(writer, "Error creaing proxy request", http.StatusInternalServerError)
		return
	}

	for name, values := range req.Header {
		for _, value := range values {
			proxyReq.Header.Add(name, value)
		}
	}

	customTransport := http.DefaultTransport
	resp, err := customTransport.RoundTrip(proxyReq)
	if err != nil {
		fmt.Println(err)
		http.Error(writer, "Error sending proxy request", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	for name, values := range resp.Header {
		for _, value := range values {
			writer.Header().Add(name, value)
		}
	}

	writer.WriteHeader(resp.StatusCode)

	io.Copy(writer, resp.Body)

	finishReq := time.Since(startReq).Seconds()
	fmt.Printf("Request performs  %s\n", finishReq)
	responseTimeHistogram.Observe(finishReq)
	requestsTotal.Inc()
}