package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	// /home/a4lab2/Documents/Projects/thoughtbin/pkg/models/Sqlite/thoughts.go
	"github.com/golangcollege/sessions"
	"gorm.io/driver/sqlite" // Sqlite driver based on GGO

	// "github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"a4lab2.com/thoughtbin/pkg/models"
	"a4lab2.com/thoughtbin/pkg/models/sq"
	"gorm.io/gorm"
)

type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

type Config struct {
	Addr      string
	StaticDir string
}

type Application struct {
	errorlog *log.Logger
	infoLog  *log.Logger
	session  *sessions.Session

	//Make the models an interface whereby anything that satisfy thier conds will be accepted, to enable us use it in testing
	thoughts interface {
		Insert(string, string, string) (uint, error)
		Get(uint) (*models.Thought, error)
		Latest() ([]*models.Thought, error)
	}
	users interface {
		Insert(string, string, string) error
		Authenticate(string, string) (uint, error)
		Get(uint) (*models.User, error)
	}
	templateCache map[string]*template.Template
}

type Thoughts struct {
}

func main() {
	var cfg Config
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")

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
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	session.SameSite = http.SameSiteStrictMode

	app := &Application{
		errorlog: errorLog,
		infoLog:  infoLog,
		thoughts: &sq.ThoughtModel{
			DB: db,
		},
		users: &sq.UserModel{
			DB: db,
		},
		session:       session,
		templateCache: templateCache,
	}
	AutoMigrate(db)
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}
	srv := &http.Server{

		Addr:      cfg.Addr,
		ErrorLog:  errorLog,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
		//timeouts for requests
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Listening on port %s", cfg.Addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	if err != nil {
		errorLog.Fatal(err)
	}
}

func AutoMigrate(conn *gorm.DB) {
	conn.Debug().AutoMigrate(
		&models.Thought{},
		&models.User{},
	)
}
