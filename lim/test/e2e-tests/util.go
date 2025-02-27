package e2eTests

import (
	"os"
)

func IsTestsFailed() bool {
	return (os.Getenv("INTEGRATION_SUCCESS") != "1") || (os.Getenv("UNIT_SUCCESS") != "1")
}
