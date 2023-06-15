package ServiceCore

import (
	"github.com/rs/zerolog/log"
	"net/rpc"
	"time"
)

func (galaxy *GalaxyClient) Call(method string, args any, reply any) error {
	serviceUrl, err := galaxy.LookUp(method)

	log.Info().Str("method", method).Str("url", serviceUrl).Msg("Looking up service for method")

	if err != nil {
		log.Err(err).Str("method", method).Msg("Unable to find a service for method")
	}

	client, err := rpc.DialHTTP("tcp", serviceUrl)

	if err != nil {
		log.Err(err).Str("method", method).Msg("Unable to call service for method")
	}

	start := time.Now()
	err = client.Call(method, args, reply)
	elapsed := time.Since(start)

	log.Info().Str("method", method).Dur("elapsed", elapsed).Msg("Called service for method: " + method)

	if err != nil {
		log.Err(err).Str("method", method).Msg("Failed when calling service for method")
	}

	return err
}