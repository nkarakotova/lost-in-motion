// +build integration

package tests

import (
	"fmt"
	"os"
	"testing"

	"lim/internal/lim-repo/flags"
	core_tests "lim/test/integration-core-tests"
	repo_tests "lim/test/integration-repo-tests"

	"github.com/charmbracelet/log"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func BeforeAllIntegration() {
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

func TestIntegrationRunner(t *testing.T) {
	BeforeAllIntegration()

	suite.RunSuite(t, new(repo_tests.ClientIntegrationRepoSuite))
	suite.RunSuite(t, new(repo_tests.TrainingIntegrationRepoSuite))
	suite.RunSuite(t, new(repo_tests.HallIntegrationRepoSuite))
	suite.RunSuite(t, new(core_tests.ClientIntegrationCoreSuite))
	suite.RunSuite(t, new(core_tests.TrainingIntegrationCoreSuite))
	suite.RunSuite(t, new(core_tests.CoachIntegrationCoreSuite))
	suite.RunSuite(t, new(core_tests.HallIntegrationCoreSuite))
}
