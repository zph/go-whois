package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/zph/go-whois/whois"
	"strings"
)

func main() {

	m := martini.Classic()
	m.Use(render.Renderer())

	m.Get("/", func(params martini.Params) string {
		return "Serving up jwhois data at /:domain"
	})

	m.Get("/whois/upload/csv", func(r render.Render) {
		r.HTML(200, "index", "foo")
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
