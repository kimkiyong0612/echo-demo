package main

import (
	"echo-demo/api/model"
	"echo-demo/api/web"
	"fmt"
	"net/url"
	"os"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func main() {

	// connect to DB
	db := sqlx.MustConnect("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true&columnsWithAlias=true&loc=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		os.Getenv("DB_NAME"),
		"Asia%2FTokyo",
	))

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
	api, err := web.NewAPI(repo, appURL)
	if err != nil {
		log.Error("[ERROR] NewAPI:", err.Error())
		return
	}

	// Create router
	e := NewRouter(api)
	// run server
	e.Logger.Fatal(e.Start(os.Getenv("PORT")))

}
