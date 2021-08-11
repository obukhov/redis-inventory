package server

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog"
	"net/http"
	"os"
)

type ServerInterface interface {
	Serve(port int, content string)
}

func NewServer(logger zerolog.Logger) *Server {
	return &Server{
		logger: logger,
	}
}

type Server struct {
	content string
	logger  zerolog.Logger
}

func (s *Server) Serve(port int, content string) {
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

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	countBytes, err := w.Write([]byte(s.content))
	if err != nil {
		s.logger.Fatal().Err(err).Msg("Error handling request")
	} else {
		s.logger.Debug().Msgf("Handled request with %d bytes", countBytes)
	}
}
