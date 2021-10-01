package main

import (
	"echo-demo/api/model"
	"echo-demo/api/web"
	"flag"
	"fmt"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"

	"github.com/gavv/httpexpect/v2"
)

var (
	testDB  string
	testSQL *sqlx.DB
	testAPI *web.API
)

var (
	flagKeepDB = flag.Bool("test-keepdb", false, "if true, don't drop test DB")
)

func TestMain(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	// load env file
	if os.Getenv("GO_ENV") == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	} else {
		if err := godotenv.Load(fmt.Sprintf(".env.%s", os.Getenv("GO_ENV"))); err != nil {
			log.Fatal(fmt.Sprintf("Error loading .env.%s file", os.Getenv("GO_ENV")))
		}
	}

	testDB = fmt.Sprintf("test_%d", time.Now().UTC().Unix())

	// DB connect
	dbRoot := sqlx.MustConnect("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s)/?parseTime=true&columnsWithAlias=true&loc=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		"Asia%2FTokyo",
	))

	// create test table
	dbRoot.MustExec("CREATE DATABASE " + testDB)
	defer func() {
		if !*flagKeepDB {
			dbRoot.MustExec("DROP DATABASE IF EXISTS " + testDB)
		}
	}()

	db := sqlx.MustConnect("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true&columnsWithAlias=true&loc=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		testDB,
		"Asia%2FTokyo",
	))

	// migrate
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}
	if _, err := migrate.Exec(db.DB, "mysql", migrations, migrate.Up); err != nil {
		log.Println("[ERROR] Failed to migrate:", err)
		return
	}

	// create repository
	repo, err := model.NewSqlxRepository(db)
	if err != nil {
		log.Println("[ERROR] NewSqlxRepository:", err.Error())
		return
	}

	appURL, err := url.Parse(os.Getenv("APP_URL"))
	if err != nil {
		log.Println("[ERROR] url.Parse(os.Getenv(\"APP_URL\"):", err.Error())
		return
	}

	// create handler
	testApi, err := web.NewAPI(repo, appURL)
	if err != nil {
		log.Error("[ERROR] NewAPI:", err.Error())
		return
	}

	// create router
	router := NewRouter(testApi)

	// run test server
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  server.URL,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	// run test
	UserTest(e)
}
