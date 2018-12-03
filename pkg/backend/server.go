package backend

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// ...
type Server struct {
	Router *mux.Router
	DB     *Database
}

// ...
func NewServer(staticPath string, connection string) (*Server, error) {
	// Set up router
	router := mux.NewRouter()

	// Set up domains
	db, err := NewDatabase(connection)
	if err != nil {
		return nil, err
	}

	// Set up server
	m := &Server{
		Router: router,
		DB:     db,
	}

	// API endpoints
	router.Methods(http.MethodGet).Path("/api/domains/").HandlerFunc(m.handleDomainsAll())

	// Debug endpoints
	router.Methods(http.MethodGet).Path("/debug/health/").HandlerFunc(m.handleDebugHealth())

	// Static catch-all
	router.Methods(http.MethodGet).PathPrefix("/").Handler(http.FileServer(http.Dir(staticPath)))

	return m, nil
}

func (m *Server) Close() {
	m.DB.Close()
}

// ...
func (m *Server) handleDomainsAll() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Write response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(m.DB.Domains())
	}
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
