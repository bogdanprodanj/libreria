// +build integration

package integration

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/libreria/config/reader"
	"github.com/libreria/models"
	"github.com/libreria/storage/postgres"
	"github.com/libreria/storage/postgres/migrations"
	"github.com/stretchr/testify/suite"
)

type LibreriaTestSuite struct {
	suite.Suite
	db *pg.DB
	c  *http.Client
}

type testConfig struct {
	PostgresTest postgres.Config `mapstructure:"postgres_test"`
}

func (s *LibreriaTestSuite) SetupSuite() {
	cfg := new(testConfig)
	err := reader.Read(cfg)
	if err != nil {
		s.FailNow("error reading configuration", err)
	}
	dir, _ := os.Getwd()
	var migD string
	if path.Base(dir) == "libreria" {
		migD = path.Join(dir, "storage", "postgres")
	} else {
		s := strings.Split(dir, "libreria")
		migD = path.Join(s[0], "libreria", "storage", "postgres")
	}
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgresTest.User, cfg.PostgresTest.Password, cfg.PostgresTest.Host, cfg.PostgresTest.Port, cfg.PostgresTest.Name)
	err = migrations.Migrate(connString, "down", migD)
	if err != nil {
		s.FailNow("failed to run migrations down", err)
	}
	err = migrations.Migrate(connString, "up", migD)
	if err != nil {
		s.FailNow("failed to run migrations up", err)
	}
	s.db = pg.Connect(&pg.Options{
		Addr:         cfg.PostgresTest.Host + ":" + cfg.PostgresTest.Port,
		User:         cfg.PostgresTest.User,
		Password:     cfg.PostgresTest.Password,
		Database:     cfg.PostgresTest.Name,
		WriteTimeout: cfg.PostgresTest.WriteTimeout,
		ReadTimeout:  cfg.PostgresTest.ReadTimeout,
		MaxRetries:   cfg.PostgresTest.MaxRetries,
	})
	s.c = &http.Client{Timeout: time.Second * 60}
}

func (s *LibreriaTestSuite) TearDownSuite() {
	_ = s.db.Close()
}

var testBooks = []*models.Book{
	{
		Title:       "Behave: The Biology of Humans at Our Best and Worst",
		Author:      "Robert M. Sapolsky",
		Publisher:   "Penguin Press",
		PublishDate: time.Date(2017, time.May, 2, 0, 0, 0, 0, time.UTC),
	},
	{
		Title:       "Can You Make This Thing Go Faster?",
		Author:      "Jeremy Clarkson",
		Publisher:   "Penguin",
		PublishDate: time.Date(2020, time.October, 29, 0, 0, 0, 0, time.UTC),
	},
	{
		Title:       "12 Rules for Life: An Antidote to Chaos",
		Author:      "Jordan B. Peterson",
		Publisher:   "Penguin Allen Lane",
		PublishDate: time.Date(2018, time.January, 16, 0, 0, 0, 0, time.UTC),
	},
}

func (s *LibreriaTestSuite) SetupTest() {
	_, err := s.db.Exec("TRUNCATE books RESTART IDENTITY")
	if err != nil {
		s.Fail("failed to truncate books table", err)
	}
	_, err = s.db.Model(&testBooks).Insert()
	if err != nil {
		s.Fail("failed to insert test values", err)
	}
}

func TestLibreriaTestSuite(t *testing.T) {
	suite.Run(t, new(LibreriaTestSuite))
}
