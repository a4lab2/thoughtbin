package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	// /home/a4lab2/Documents/Projects/thoughtbin/pkg/models/Sqlite/thoughts.go
	"gorm.io/driver/sqlite" // Sqlite driver based on GGO
	// "github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"a4lab2.com/thoughtbin/pkg/models"
	"a4lab2.com/thoughtbin/pkg/models/sq"
	"gorm.io/gorm"
)

type Config struct {
	Addr      string
	StaticDir string
}

type Application struct {
	errorlog *log.Logger
	infoLog  *log.Logger
	thoughts *sq.ThoughtModel
}

type Thoughts struct {
}

func main() {
	var cfg Config
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()
	//Logs
	infoLog := log.New(os.Stdout, "INFOLOG\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ErrorLog\t", log.Ldate|log.Ltime|log.Lshortfile)

	// dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	// flag.Parse()

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		errorLog.Fatal(err)
	}
	sqdb, err := db.DB()
	if err != nil {
		errorLog.Fatal(err)
	}
	// sqdb.Ping()
	defer sqdb.Close()

	app := &Application{
		errorlog: errorLog,
		infoLog:  infoLog,
		thoughts: &sq.ThoughtModel{
			DB: db,
		},
	}
	AutoMigrate(db)
	srv := &http.Server{

		Addr:     cfg.Addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Listening on port %s", cfg.Addr)
	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}

func AutoMigrate(conn *gorm.DB) {
	conn.Debug().AutoMigrate(
		&models.Thought{},
	)
}
