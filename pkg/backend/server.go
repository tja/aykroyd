package backend

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/tja/postfix-web/pkg/assets"
)

// Server holds the state of a HTTP server. The HTTP server exposes REST methods to manipulate domains and
// their associated (email) forwards.
type Server struct {
	Router *mux.Router
	DB     *Database
}

// NewServer creates a new server. assetPath points to a folder with static web content. If assetPath is empty,
// the web content is loaded from the embedded filesystem. mysql holds the MySQL connection string, which is
// passed through to the database layer.
func NewServer(assetPath string, mysql string) (*Server, error) {
	// Set up router
	router := mux.NewRouter()

	// Set up domains
	db, err := NewDatabase(mysql)
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
	if assetPath == "" {
		// Load from embedded filesystem
		logrus.Debug("Serving embedded assets")
		router.Methods(http.MethodGet).PathPrefix("/").Handler(http.FileServer(assets.HTTP))
	} else {
		// Load from given path
		logrus.Debugf("Serving assets from path '%s'", assetPath)
		router.Methods(http.MethodGet).PathPrefix("/").Handler(http.FileServer(http.Dir(assetPath)))
	}

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
		logrus.Info("Return domains")

		// Get domains
		domains, err := m.DB.Domains()
		if err != nil {
			// Unable to fetch domains
			w.WriteHeader(http.StatusInternalServerError)

			logrus.Warn(err)
			return
		}

		// Convert to Json
		json, err := json.Marshal(domains)
		if err != nil {
			// Unable to generate Json
			w.WriteHeader(http.StatusInternalServerError)

			logrus.Warn(err)
			return
		}

		// Write success response
		w.WriteHeader(http.StatusOK)
		w.Write(json)

		logrus.Debugf("Returned %d domains", len(domains))
	}
}

// handleDomainsCreateDomain returns a function that handles incoming REST requests asking to create a new
// domain.
func (m *Server) handleDomainsCreateDomain() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logrus.Info("Create new domain")

		// Parse input
		var input struct {
			Name string `json:"name"`
		}

		err := json.NewDecoder(req.Body).Decode(&input)
		if err != nil {
			// Could not be parsed
			w.WriteHeader(http.StatusBadRequest)

			logrus.Warn(err)
			return
		}

		// Create domain in database
		err = m.DB.CreateDomain(input.Name)
		if err != nil {
			// Domain already exists
			w.WriteHeader(http.StatusInternalServerError)

			logrus.Warn(err)
			return
		}

		// Write success response
		w.WriteHeader(http.StatusOK)

		logrus.Debugf("Created new domain '%s'", input.Name)
	}
}

// handleDomainsDeleteDomain returns a function that handles incoming REST requests asking to delete an
// existing domain. All associated forwards are deleted as well.
func (m *Server) handleDomainsDeleteDomain() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logrus.Info("Delete domain")

		// Extract URI variables
		domain := mux.Vars(req)["domain"]

		// Create forward in database
		err := m.DB.DeleteDomain(domain)
		if err != nil {
			// Forward already exists
			w.WriteHeader(http.StatusInternalServerError)

			logrus.Warn(err)
			return
		}

		// Write success response
		w.WriteHeader(http.StatusOK)

		logrus.Debugf("Deleted domain '%s'", domain)
	}
}

// handleDomainsCreateForward returns a function that handles incoming REST requests asking to create a new
// forward for an existing domain.
func (m *Server) handleDomainsCreateForward() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logrus.Info("Create email forward")

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

			logrus.Warn(err)
			return
		}

		// Create forward in database
		err = m.DB.CreateForward(domain, input.From, input.To)
		if err != nil {
			// Forward already exists
			w.WriteHeader(http.StatusInternalServerError)

			logrus.Warn(err)
			return
		}

		// Write success response
		w.WriteHeader(http.StatusOK)

		logrus.Debugf("Created email forward for domain '%s' from '%s' to '%s'", domain, input.From, input.To)
	}
}

// handleDomainsUpdateForward returns a function that handles incoming REST requests asking to update an
// existing forward of an existing domain.
func (m *Server) handleDomainsUpdateForward() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logrus.Info("Update email forward")

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

			logrus.Warn(err)
			return
		}

		// Create forward in database
		err = m.DB.UpdateForward(domain, from, input.To)
		if err != nil {
			// Forward already exists
			w.WriteHeader(http.StatusInternalServerError)

			logrus.Warn(err)
			return
		}

		// Write success response
		w.WriteHeader(http.StatusOK)

		logrus.Debugf("Updated email forward '%s' of domain '%s' to '%s'", from, domain, input.To)
	}
}

// handleDomainsDeleteForward returns a function that handles incoming REST requests asking to delete an
// existing forward of an existing domain.
func (m *Server) handleDomainsDeleteForward() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logrus.Info("Delete email forward")

		// Extract URI variables
		vars := mux.Vars(req)

		domain := vars["domain"]
		from := vars["from"]

		// Create forward in database
		err := m.DB.DeleteForward(domain, from)
		if err != nil {
			// Forward already exists
			w.WriteHeader(http.StatusInternalServerError)

			logrus.Warn(err)
			return
		}

		// Write success response
		w.WriteHeader(http.StatusOK)

		logrus.Debugf("Deleted email forward '%s' of domain '%s'", from, domain)
	}
}
