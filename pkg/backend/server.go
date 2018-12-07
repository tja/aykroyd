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
	router.Methods(http.MethodDelete).Path("/api/domains/{domain}/").HandlerFunc(m.handleDomainsDeleteDomain())

	router.Methods(http.MethodPost).Path("/api/domains/{domain}/forwards/").HandlerFunc(m.handleDomainsCreateForward())
	router.Methods(http.MethodPut).Path("/api/domains/{domain}/forwards/{from}/").HandlerFunc(m.handleDomainsChangeForward())
	router.Methods(http.MethodDelete).Path("/api/domains/{domain}/forwards/{from}/").HandlerFunc(m.handleDomainsDeleteForward())

	// Static catch-all
	router.Methods(http.MethodGet).PathPrefix("/").Handler(http.FileServer(http.Dir(staticPath)))

	return m, nil
}

// ...
func (m *Server) Close() error {
	return m.DB.Close()
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
func (m *Server) handleDomainsDeleteDomain() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Extract domain
		domain := mux.Vars(req)["domain"]

		// Create forward in database
		err := m.DB.DeleteDomain(domain)
		if err != nil {
			// Forward already exists
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Write response
		w.WriteHeader(http.StatusOK)
	}
}

// ...
func (m *Server) handleDomainsCreateForward() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Extract domain
		domain := mux.Vars(req)["domain"]

		// Parse input
		var input struct {
			From string `json:"from"`
			To   string `json:"to"`
		}

		err := json.NewDecoder(req.Body).Decode(&input)
		if err != nil {
			// Could not be parsed
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Create forward in database
		err = m.DB.CreateForward(domain, input.From, input.To)
		if err != nil {
			// Forward already exists
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Write response
		w.WriteHeader(http.StatusOK)
	}
}

// ...
func (m *Server) handleDomainsChangeForward() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Extract URL variables
		vars := mux.Vars(req)

		domain := vars["domain"]
		from := vars["from"]

		// Parse input
		var input struct {
			To string `json:"to"`
		}

		err := json.NewDecoder(req.Body).Decode(&input)
		if err != nil {
			// Could not be parsed
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Create forward in database
		err = m.DB.UpdateForward(domain, from, input.To)
		if err != nil {
			// Forward already exists
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Write response
		w.WriteHeader(http.StatusOK)
	}
}

// ...
func (m *Server) handleDomainsDeleteForward() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Extract URL variables
		vars := mux.Vars(req)

		domain := vars["domain"]
		from := vars["from"]

		// Create forward in database
		err := m.DB.DeleteForward(domain, from)
		if err != nil {
			// Forward already exists
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Write response
		w.WriteHeader(http.StatusOK)
	}
}
