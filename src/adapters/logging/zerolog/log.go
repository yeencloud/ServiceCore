package zerolog

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yeencloud/ServiceCore/src/domain"
	"github.com/yeencloud/ServiceCore/src/domain/types"
	"os"
	"strings"
)

type Logger struct {
}

func (l Logger) Debug(message string, fields ...domain.LogField) {
	//TODO implement me
	panic("implement me")
}

func (l Logger) Info(message string, fields ...domain.LogField) {
	//TODO implement me
	panic("implement me")
}

func (l Logger) Warn(message string, fields ...domain.LogField) {
	//TODO implement me
	panic("implement me")
}

func (l Logger) Error(message string, fields ...domain.LogField) {
	//TODO implement me
	panic("implement me")
}

func NewLogger(env types.Environment) Logger {
	isDev := env.IsDevelopment()

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
	return Logger{}
}
