package backend

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// ...
type Server struct {
	Router *mux.Router
}

// ...
func NewServer(staticPath string) (*Server, error) {
	// Set up router
	router := mux.NewRouter()

	m := &Server{
		Router: router,
	}

	// Set up API endpoints
	router.Methods(http.MethodGet).Path("/debug/health/").HandlerFunc(m.handleDebugHealth())

	// Static catch-all
	router.Methods(http.MethodGet).PathPrefix("/").Handler(http.FileServer(http.Dir(staticPath)))

	return m, nil
}

// ...
func (m *Server) handleDebugHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// ...
		body := map[string]string{
			"response": "healthy",
		}

		// Write response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(body)
	}
}
