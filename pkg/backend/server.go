package backend

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Server holds the state of a HTTP server. The HTTP server exposes REST methods to manipulate domains and
// their associated (email) forwards.
type Server struct {
	Router *mux.Router
	DB     *Database
}

// NewServer creates a new server. staticPath points to a folder with all static assets. connection holds the
// MySQL connection string, which is passed through to the database layer.
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

	// REST methods
	router.Methods(http.MethodGet).Path("/api/domains/").HandlerFunc(m.handleDomainsListDomains())

	router.Methods(http.MethodPost).Path("/api/domains/").HandlerFunc(m.handleDomainsCreateDomain())
	router.Methods(http.MethodDelete).Path("/api/domains/{domain}/").HandlerFunc(m.handleDomainsDeleteDomain())

	router.Methods(http.MethodPost).Path("/api/domains/{domain}/forwards/").HandlerFunc(m.handleDomainsCreateForward())
	router.Methods(http.MethodPut).Path("/api/domains/{domain}/forwards/{from}/").HandlerFunc(m.handleDomainsUpdateForward())
	router.Methods(http.MethodDelete).Path("/api/domains/{domain}/forwards/{from}/").HandlerFunc(m.handleDomainsDeleteForward())

	// Static catch-all
	router.Methods(http.MethodGet).PathPrefix("/").Handler(http.FileServer(http.Dir(staticPath)))

	return m, nil
}

// Close closes the server.
func (m *Server) Close() error {
	return m.DB.Close()
}

// handleDomainsListDomains returns a function that handles incoming REST requests asking for the list of all
// domains and associated forwards.
func (m *Server) handleDomainsListDomains() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Get domains
		domains, err := m.DB.Domains()
		if err != nil {
			// Unable to fetch domains
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Convert to Json
		json, err := json.Marshal(domains)
		if err != nil {
			// Unable to generate Json
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Write success response
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}

// handleDomainsCreateDomain returns a function that handles incoming REST requests asking to create a new
// domain.
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

		// Write success response
		w.WriteHeader(http.StatusOK)
	}
}

// handleDomainsDeleteDomain returns a function that handles incoming REST requests asking to delete an
// existing domain. All associated forwards are deleted as well.
func (m *Server) handleDomainsDeleteDomain() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Extract URI variables
		domain := mux.Vars(req)["domain"]

		// Create forward in database
		err := m.DB.DeleteDomain(domain)
		if err != nil {
			// Forward already exists
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Write success response
		w.WriteHeader(http.StatusOK)
	}
}

// handleDomainsCreateForward returns a function that handles incoming REST requests asking to create a new
// forward for an existing domain.
func (m *Server) handleDomainsCreateForward() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Extract URI variables
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

		// Write success response
		w.WriteHeader(http.StatusOK)
	}
}

// handleDomainsUpdateForward returns a function that handles incoming REST requests asking to update an
// existing forward of an existing domain.
func (m *Server) handleDomainsUpdateForward() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Extract URI variables
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

		// Write success response
		w.WriteHeader(http.StatusOK)
	}
}

// handleDomainsDeleteForward returns a function that handles incoming REST requests asking to delete an
// existing forward of an existing domain.
func (m *Server) handleDomainsDeleteForward() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Extract URI variables
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

		// Write success response
		w.WriteHeader(http.StatusOK)
	}
}
