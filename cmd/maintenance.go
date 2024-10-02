package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"treetime/pkg/utils/logging"
	"treetime/pkg/utils/net"
)

func maintenanceMode() {

	var port = os.Getenv("PORT")

	var maintenancePageTemplate *template.Template
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		if !net.IsAjax(r) {
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