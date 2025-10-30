package main

import (
	"log"

	"github.com/go-fuego/fuego"
	"github.com/mesbrj/dbus-controller/internal/api"
	"github.com/mesbrj/dbus-controller/internal/service"
)

func main() {
	// Create D-Bus service
	dbusService := service.NewDBusService()
	defer dbusService.Close()

	// Create Fuego server
	s := fuego.NewServer(
		fuego.WithAddr(":8080"),
	)

	// Setup routes
	api.SetupRoutes(s, dbusService)

	// Start server
	log.Println("Starting D-Bus Controller API on :8080")
	if err := s.Run(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
