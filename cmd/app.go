package cmd

import (
	"caching-proxy/internal/cache"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

// ClearCache clears all the data in the cache
func ClearCache() {
	log.Println("Clearing caches...")
	cache.CacheManager.DeleteAll()
	log.Println("Cache cleared successfully")
}

func StartServer(port, origin string) {
	// Handle request
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleRequest(w, r, origin)
	})

	// Start the server
	log.Printf("Starting server on port %s with origin %s\n", port, origin)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request, origin string) {
	// Get the endpoint from the request and check if it exists in the cache
	endpoint := r.URL.Path
	cacheData, ok := cache.CacheManager.Get(endpoint)
	if ok {
		// If the data exists in the cache, check if it's expired
		if cacheData.ExpiredAt.After(cacheData.CreatedAt) {
			// If the data is not expired, return it
			log.Printf("Cache hit for endpoint: %s\n", endpoint)
			w.Header().Set("X-Cache", "HIT")
			_, err := w.Write(cacheData.Data)
			if err != nil {
				return
			}
			return
		}
	}

	// If the data is not in the cache or expired, fetch it from the origin server
	log.Printf("Cache miss for endpoint: %s\n", endpoint)

	// Join the origin and endpoint to get the full path
	path, err := url.JoinPath(origin, endpoint)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			return
		}
		return
	}

	// Make a request to the origin server
	resp, err := http.Get(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			return
		}
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	// Read the response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			return
		}
		return
	}

	// Set the data in the cache
	cacheData = cache.CacheData{
		Data:      data,
		Endpoint:  endpoint,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(time.Duration(60) * time.Second),
	}
	cache.CacheManager.Set(endpoint, cacheData)

	// Write the data to the response writer
	w.Header().Set("X-Cache", "MISS")
	_, err = w.Write(data)
	if err != nil {
		return
	}
	return
}
