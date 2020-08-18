package ldap

import (
	"fmt"
	"github.com/nmcclain/ldap"
	"github.com/radovskyb/watcher"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	query "github.com/shauncampbell/dapper/pkg/query"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// Server is a struct holding all of the state information about your LDAP server
type Server struct {
	port       int            // The port number of the LDAP server
	baseDN     string         // The baseDN that the LDAP server will service
	configFile string         // The configFile containing the config for the server
	s          *ldap.Server   // The underlying ldap.Server implementation
	logger     zerolog.Logger // The logger being used for console printing
	lock       sync.Mutex     // A lock to prevent multiple updates clashing
	entries    []*ldap.Entry  // The ldap entries read from the configuration file
}

// NewServer creates a new server instance which manages a given baseDN and stores
// user information in the specified configFile.
func NewServer(baseDN, configFile string, port int) *Server {
	server := &Server{baseDN: baseDN, configFile: configFile, port: port, logger: log.Output(zerolog.ConsoleWriter{Out: os.Stderr})}
	s := ldap.NewServer()
	server.s = s

	// register Bind and Search function handlers
	s.BindFunc(baseDN, server)
	s.SearchFunc(baseDN, server)

	return server
}

// Listen starts the server up on the specified port and begins listening for connections.
func (s *Server) Listen() error {
	// start the server
	listen := fmt.Sprintf("0.0.0.0:%d", s.port)
	s.logger.Info().Msgf("starting LDAP server on %s for %s", listen, s.baseDN)

	// Start waiting for configuration changes
	go s.WatchForConfigChanges()

	// Load the initial configuration the first time
	s.ReloadConfiguration(s.configFile)

	// Listen
	return s.s.ListenAndServe(listen)
}

// WatchForConfigChanges starts watching the configuration file for writes and applies changes automatically.
func (s *Server) WatchForConfigChanges() {
	// Set up the file watcher
	w := watcher.New()

	// SetMaxEvents to 1 to allow at most 1 event's to be received
	// on the Event channel per watching cycle.
	//
	// If SetMaxEvents is not set, the default is to send all events.
	w.SetMaxEvents(1)
	go func() {
		for {
			select {
			case event := <-w.Event:
				s.ReloadConfiguration(event.Path)
			case err := <-w.Error:
				s.logger.Err(err)
			case <-w.Closed:
				return
			}
		}
	}()
	// Watch this file for changes.
	if err := w.Add(s.configFile); err != nil {
		s.logger.Error().Err(err)
		return
	}

	if err := w.Start(10 * time.Second); err != nil {
		s.logger.Error().Err(err)
		return
	}
}

// parseAttributeValues parses the attribute values from a yaml file.
func parseAttributeValues(v interface{}) ([]string, bool) {
	if str, ok := v.(string); ok {
		return []string{str}, true
	}

	if inter, ok := v.([]interface{}); ok {
		out := make([]string, 0)
		for _, val := range inter {
			if str, ok := val.(string); ok {
				out = append(out, str)
			}
		}
		return out, true
	}

	return nil, false
}

// parseUser parses the user record from the yaml file.
func parseUser(dn string, user map[interface{}]interface{}, logger zerolog.Logger) (*ldap.Entry, error) {
	entry := &ldap.Entry{
		DN:         dn,
		Attributes: make([]*ldap.EntryAttribute, 0),
	}

	for k, v := range user {
		// skip the dn key as we already know about that one.
		if k == dn {
			continue
		}

		name, ok := k.(string)
		if !ok {
			logger.Warn().Msgf("skipping attribute dn '%s' because there is no valid name", dn)
			continue
		}

		value, ok := parseAttributeValues(v)
		if !ok {
			logger.Warn().Msgf("skipping attribute '%s' for dn '%s' because there is no valid values", name, dn)
			continue
		}

		// apply transformation for password
		if name == "userPassword" {
			for i, v := range value {
				p, err := parsePassword(v)
				if err != nil {
					logger.Error().Err(err).Msg("failed to parse password")
					continue
				}
				value[i] = p
			}
		}

		logger.Info().Str("dn", dn).Str("attribute", name).Strs("value", value).Msgf("adding attribute")

		entry.Attributes = append(entry.Attributes, &ldap.EntryAttribute{Name: name, Values: value})
	}
	return entry, nil
}

// parseUsers parses all user records from the yaml file.
func parseUsers(users []interface{}, logger zerolog.Logger) ([]*ldap.Entry, error) {
	entries := make([]*ldap.Entry, 0)
	for _, user := range users {
		if u, ok := user.(map[interface{}]interface{}); ok {
			if dn, ok := u["dn"].(string); ok {
				entry, err := parseUser(dn, u, logger)
				if err == nil {
					entries = append(entries, entry)
				} else {
					logger.Warn().Interface("user", u).Err(err).Msg("user record invalid")
					return nil, err
				}
			}
		}
	}
	return entries, nil
}

// ReloadConfiguration reads the configuration file and applies the changes.
func (s *Server) ReloadConfiguration(filename string) {
	s.logger.Info().Msgf("reloading configuration file '%s'", filename)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to read file")
		return
	}
	var c map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to unmarshall file")
		return
	}

	// extract users
	if c["users"] != nil {
		var entries []*ldap.Entry
		if users, ok := c["users"].([]interface{}); ok {
			entries, err = parseUsers(users, s.logger)
			if err != nil {
				s.logger.Error().Err(err).Msg("failed to parse users")
				return
			}
		}

		if entries != nil {
			s.lock.Lock()
			s.entries = entries
			defer s.lock.Unlock()
		}
	}
}

// parsePassword parses a password string and ensures that it is in SSHA format.
func parsePassword(password string) (string, error) {
	if strings.HasPrefix(password, "{SSHA}") {
		return password, nil
	} else if strings.HasPrefix(password, "{") {
		return "", fmt.Errorf("unsupported password encoding scheme")
	} else {
		encoder := SSHAEncoder{}
		ssha, err := encoder.Encode([]byte(password))
		if err != nil {
			return "", err
		} else {
			return string(ssha), nil
		}
	}
}

// Bind is a handler for an incoming bind request.
func (s *Server) Bind(bindDN, bindSimplePw string, conn net.Conn) (ldap.LDAPResultCode, error) {
	logger := s.logger.With().Str("operation", "bind").Str("request_ip", conn.RemoteAddr().String()).Str("bindDN", bindDN).Logger()
	logger.Info().Msgf("request received")
	for _, entry := range s.entries {
		if entry.DN == bindDN {
			encoder := SSHAEncoder{}
			pwd := entry.GetAttributeValue("userPassword")
			if encoder.Matches([]byte(pwd), []byte(bindSimplePw)) {
				logger.Info().Msgf("bind request was accepted")
				return ldap.LDAPResultSuccess, nil
			}
			logger.Error().Msgf("bind request was rejected because of an invalid password")
			return ldap.LDAPResultInvalidCredentials, nil
		}
	}

	logger.Error().Msgf("bind request was rejected because the dn does not exist")
	return ldap.LDAPResultInvalidCredentials, nil
}

// Search is a handler for an incoming search request.
func (s *Server) Search(boundDN string, searchReq ldap.SearchRequest, conn net.Conn) (ldap.ServerSearchResult, error) {
	logger := s.logger.With().Str("operation", "search").Str("request_ip", conn.RemoteAddr().String()).Str("bindDN", boundDN).Logger()

	// Parse the search query
	logger.Info().Msgf("beginning search with query: %s", searchReq.Filter)
	q, _, err := query.Parse(searchReq.Filter, 0)
	if err != nil {
		logger.Error().Err(err).Msgf("the client submitted an invalid query: %s", searchReq.Filter)
		return ldap.ServerSearchResult{Entries: nil, Referrals: nil, Controls: nil, ResultCode: ldap.LDAPResultUnwillingToPerform}, nil
	}

	var result = make([]*ldap.Entry, 0)
	for _, entry := range s.entries {
		if q.Evaluate(entry) {
			logger.Info().Msgf("dn '%s' matches search criteria", entry.DN)
			result = append(result, entry)
		}
	}
	logger.Info().Msgf("search completed with %d results", len(result))

	return ldap.ServerSearchResult{Entries: result, Referrals: []string{}, Controls: []ldap.Control{}, ResultCode: ldap.LDAPResultSuccess}, nil
}
