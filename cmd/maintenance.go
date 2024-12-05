package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"spacetime/pkg/utils/logging"
	"spacetime/pkg/utils/net"
)

func maintenanceMode() {

	var port = os.Getenv("PORT")

	var maintenancePageTemplate *template.Template

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		if net.IsAjax(r) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error": true, "maintenanceMode": true}`))
		} else {
			if maintenancePageTemplate == nil {
				maintenancePageTemplate = template.Must(template.ParseFiles("html/maintenance.html"))
			}
			maintenancePageTemplate.Execute(w, nil)
		}
	})

	log.Printf("Listening on port %s", port)

	if err := http.ListenAndServe(":"+port, nil); err != http.ErrServerClosed {
		logging.LogErrorFatal(err)
	}

	// Flush pending logs
	logging.CloseLoggingClients()

}
