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
	router.Methods(http.MethodGet).Path("/api/domains/").HandlerFunc(m.handleDomainsListDomains())
	router.Methods(http.MethodPost).Path("/api/domains/").HandlerFunc(m.handleDomainsCreateDomain())

	// Debug endpoints
	router.Methods(http.MethodGet).Path("/debug/health/").HandlerFunc(m.handleDebugHealth())

	// Static catch-all
	router.Methods(http.MethodGet).PathPrefix("/").Handler(http.FileServer(http.Dir(staticPath)))

	return m, nil
}

// ...
func (m *Server) Close() {
	m.DB.Close()
}

// ...
func (m *Server) handleDomainsListDomains() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Get domains as Json
		json, err := json.Marshal(m.DB.ListDomains())
		if err != nil {
			// Unable to generate Json
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Write response
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}

// ...
func (m *Server) handleDomainsCreateDomain() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Parse input
		var input struct {
			Name string `json:"name"`
		}

		err := json.NewDecoder(req.Body).Decode(&input)
		if err != nil {
			// Could not be parsed
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Create domain in database
		err = m.DB.CreateDomain(input.Name)
		if err != nil {
			// Domain already exists
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Write response
		w.WriteHeader(http.StatusOK)
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
