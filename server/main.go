package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func main() {
	// Redis-Konfiguration
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis-Server-Adresse
		Password: "",               // Passwort (falls erforderlich)
		DB:       0,                // Redis-Datenbanknummer
	})

	// Überprüfe die Verbindung zur Redis-Datenbank
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Fehler beim Verbinden zur Redis-Datenbank:", err)
	}

	// Erstelle den HTTP-Handler
	http.HandleFunc("/", saveIPHandler)

	// Starte den Server auf Port 8080
	log.Println("Server gestartet. Höre auf Port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Fehler beim Starten des Servers:", err)
	}
}

func saveIPHandler(w http.ResponseWriter, r *http.Request) {
	// Ermittle die IP-Adresse der Anfrage
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Speichere die IP-Adresse in der Redis-Datenbank
	err = redisClient.SAdd(r.Context(), "ip_addresses", ip).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Gib eine Erfolgsmeldung zurück
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("IP-Adresse erfolgreich gespeichert: " + ip))
}
