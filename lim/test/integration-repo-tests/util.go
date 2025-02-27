package integrationRepoTests

import (
	"os"
)

func IsUnitTestsFailed() bool {
	return os.Getenv("UNIT_SUCCESS") != "1"
}
