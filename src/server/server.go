package server

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog"
	"net/http"
	"os"
)

// SimpleServerInterface abstraction for simplified web server
type SimpleServerInterface interface {
	Serve(port int, content string)
}

// NewServer creates SimpleServer
func NewServer(logger zerolog.Logger) *SimpleServer {
	return &SimpleServer{
		logger: logger,
	}
}

// SimpleServer is SimpleServerInterface implementation
type SimpleServer struct {
	content string
	logger  zerolog.Logger
}

// Serve initiates server, block for input on stdin
func (s *SimpleServer) Serve(port int, content string) {
	s.content = content
	go func() {
		s.logger.Info().Msgf("Listening on port %d: http://localhost:%d/", port, port)
		s.logger.Info().Msg("Press enter to exit")
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), s)
		s.logger.Fatal().Err(err).Msg("Error serving content")
	}()

	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}

// ServeHTTP handler function
func (s *SimpleServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	countBytes, err := w.Write([]byte(s.content))
	if err != nil {
		s.logger.Fatal().Err(err).Msg("Error handling request")
	} else {
		s.logger.Debug().Msgf("Handled request with %d bytes", countBytes)
	}
}
