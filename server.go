package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/zph/go-whois/whois"
	"mime/multipart"
	"strings"
	"sync"
)

type UploadForm struct {
	FileUpload *multipart.FileHeader `form:"domains"`
}

func main() {

	m := martini.Classic()
	m.Use(render.Renderer())

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
			go whois.AsyncRetrieve(line[0], messages, &wg)
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
		rec := whois.RetrieveJSON(params["domain"])
		return rec
	})

	m.Get("/whois/emails/:domain", func(params martini.Params) string {
		rec, _ := whois.Retrieve(params["domain"])
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
