package servicecore

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

// we are using zero log for logging
// it helps us logging in either json or console format depending on the environment as well as logging the file and line number for debugging purposes
func (sh *ServiceHost) setupLogging() {
	isDev := sh.Config.IsDevelopment()

	if isDev {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Short caller (file:line)
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		if strings.HasPrefix(file, "/app/") {
			file = strings.TrimPrefix(file, "/app/")
		}

		return fmt.Sprintf("%s:%d", file, line)
	}

	if isDev {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Logger = log.With().Caller().Logger()

	if isDev {
		log.Warn().Msg("Running in development mode logging can be verbose or contain sensitive information")
		log.Warn().Msg("To run in production mode set the ENV environment variable to `prod` or `production`")
	}
}