package ServiceCore

import "os"

func isDevelopment() bool {
	return !(os.Getenv("ENV") == "production" || os.Getenv("ENV") == "prod")
}