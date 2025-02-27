// +build e2e

package tests

import (
	"fmt"
	"os"
	"testing"

	"lim/internal/lim-repo/flags"
	e2e_tests "lim/test/e2e-tests"

	"github.com/charmbracelet/log"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func BeforeAllE2E() {
	pf := flags.PostgresFlags{
		Host:     "postgres",
		User:     "natali",
		Password: "12345",
		Port:     "5432",
		DBName:   "postgres",
	}

	db, err := pf.InitDB(log.Default())
	if err != nil {
		fmt.Println(err)
		return
	}

	text, err := os.ReadFile("/builds/knv21u506/test-karakotova-73/lim/database/postgreSQL/init.sql")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Exec(string(text))
	if err != nil {
		fmt.Println(err)
		return
	}

	text, err = os.ReadFile("/builds/knv21u506/test-karakotova-73/lim/database/postgreSQL/test.sql")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Exec(string(text))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestE2ERunner(t *testing.T) {
	BeforeAllE2E()

	suite.RunSuite(t, new(e2e_tests.E2ESuite))
}
