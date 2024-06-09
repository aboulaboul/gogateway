package main

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

// Route représente une route de redirection
type Route struct {
	Path   string `json:"path"`
	Target string `json:"target"`
}

// Config représente la configuration des routes
type Config struct {
	Routes []Route `json:"routes"`
}

// proxyHandler redirige les requêtes vers l'URL cible
func proxyHandler(prefix string, target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request for: %s", r.URL.Path)

		targetURL, err := url.Parse(target)
		if err != nil {
			log.Printf("Invalid target URL: %s", target)
			http.Error(w, "Invalid target URL", http.StatusInternalServerError)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
			log.Printf("Error connecting to the backend server: %v", err)
			http.Error(rw, "Error connecting to the backend server", http.StatusBadGateway)
		}

		// Ajuster le chemin de l'URL pour correspondre au service cible
		r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
		log.Printf("Rewriting URL to: %s", r.URL.Path)

		// Gérer les WebSockets
		if strings.HasPrefix(r.Header.Get("Connection"), "Upgrade") && r.Header.Get("Upgrade") == "websocket" {
			proxy.Director = func(req *http.Request) {
				req.URL.Scheme = targetURL.Scheme
				req.URL.Host = targetURL.Host
				req.URL.Path = targetURL.Path
				if targetURL.RawQuery == "" || req.URL.RawQuery == "" {
					req.URL.RawQuery = targetURL.RawQuery + req.URL.RawQuery
				} else {
					req.URL.RawQuery = targetURL.RawQuery + "&" + req.URL.RawQuery
				}
				req.Host = targetURL.Host
			}
		}

		proxy.ServeHTTP(w, r)
	}
}

func main() {
	// Lire la configuration à partir du fichier routes.json
	configFile, err := os.Open("routes.json")
	if err != nil {
		log.Fatalf("Could not open routes file: %s\n", err)
	}
	defer configFile.Close()

	byteValue, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatalf("Could not read config file: %s\n", err)
	}

	var config Config
	if err := json.Unmarshal(byteValue, &config); err != nil {
		log.Fatalf("Could not parse routes file: %s\n", err)
	}

	// Définir les routes pour les différents services
	for _, route := range config.Routes {
		log.Printf("Configuring route: %s -> %s", route.Path, route.Target)
		http.HandleFunc(route.Path, proxyHandler(route.Path, route.Target))
	}

	// Charger les certificats SSL/TLS
	certFile := "server.crt"
	keyFile := "server.key"

	// Configurer le serveur HTTPS
	server := &http.Server{
		Addr: ":443",
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	// Démarrer le serveur HTTPS
	log.Println("Starting gateway on port 443")
	if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
