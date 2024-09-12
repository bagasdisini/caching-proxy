package main

import (
	"caching-proxy/cmd"
	"flag"
	"log"
)

func main() {
	// Parse the command line flags
	port := flag.String("port", "", "Port on which the caching proxy server will run")
	origin := flag.String("origin", "", "URL of the server to which the requests will be forwarded")
	clearCache := flag.Bool("clear-cache", false, "Clear the cache")

	flag.Parse()

	// Only clear the cache if the flag is set
	if *clearCache {
		cmd.ClearCache()
		return
	}

	// Check if the required flags are provided
	if *port == "" || *origin == "" {
		log.Fatal("Please provide the port and origin URL. Example: ./caching-proxy --port 8080 --origin https://google.com")
		return
	}

	// Start the caching proxy server
	cmd.StartServer(*port, *origin)
}
