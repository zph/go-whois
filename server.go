package main

import "github.com/go-martini/martini"
import "github.com/zph/go-whois/whois"

func main() {

    m := martini.Classic()

    m.Get("/", func(params martini.Params) string {
        return "Serving up jwhois data at /:domain"
    })
    m.Get("/:domain", func(params martini.Params) string {
        rec := whois.RetrieveJSON(params["domain"])
        return rec
    })
    m.Get("/favicon.ico", func() (int) {
        return 418
    })
    m.NotFound(func() (int, string) {
        return 418, "yolo"
    })
    m.Run()
}
