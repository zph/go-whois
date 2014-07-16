package main

import (
	"database/sql"
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zph/go-whois/whois"
	"mime/multipart"
	"runtime"
	"strings"
	"sync"
)

type UploadForm struct {
	FileUpload *multipart.FileHeader `form:"domains"`
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())
	runtime.GOMAXPROCS(runtime.NumCPU())

	dbmap := init_db()
	defer dbmap.Db.Close()

	m.Get("/", func(params martini.Params) string {
		return "Serving up jwhois data at /:domain"
	})

	m.Get("/whois/upload/csv", func(r render.Render) {
		r.HTML(200, "index", "foo")
	})

	m.Post("/whois/upload/csv", binding.MultipartForm(UploadForm{}), func(uf UploadForm) string {
		file, _ := uf.FileUpload.Open()
		lines := whois.ParseCSV(file)
		count := len(lines)

		messages := make(chan string, count)

		var wg sync.WaitGroup

		for _, line := range lines {
			wg.Add(1)
			go whois.AsyncRetrieve(line[0], dbmap, messages, &wg)
		}

		go func() {
			wg.Wait()
			close(messages)
		}()

		msgs := make([]string, 0)
		for elem := range messages {
			msgs = append(msgs, elem)
		}
		return strings.Join(msgs, "\n")
	})

	m.Get("/whois/:domain", func(params martini.Params) string {
		rec := whois.RetrieveJSON(params["domain"], dbmap)
		return rec
	})

	m.Get("/whois/emails/:domain", func(params martini.Params) string {
		rec, _ := whois.Retrieve(params["domain"], dbmap)
		return strings.Join(rec.Emails, ", ")
	})

	m.Get("/favicon.ico", func() int {
		return 418
	})
	m.NotFound(func() (int, string) {
		return 418, "yolo"
	})
	m.Run()
}

func init_db() *gorp.DbMap {
	db, err := sql.Open("sqlite3", "./whois.db")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbmap.AddTableWithName(whois.SqlResult{}, "WhoisResults").SetKeys(false, "Domain", "Raw", "Emails")
	err = dbmap.CreateTablesIfNotExists()

	if err != nil {
		fmt.Println("Error initializing db")
	}
	return dbmap
}
