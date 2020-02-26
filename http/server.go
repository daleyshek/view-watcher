package http

import (
	"html/template"
	"log"
	"net/http"

	"github.com/daleyshek/view-watcher/cache"
)

// ServeAll run all
func ServeAll() {
	serveHTTP()
}

func serveHTTP() {
	http.HandleFunc("/watcher/test", index)
	http.HandleFunc("/watcher/report", handleReport)
	http.HandleFunc("/watcher/js", handleJS)
	http.ListenAndServe(":"+cache.Config.Port, nil)
}

func handleJS(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("report.js")
	if err != nil {
		log.Fatal("report.js was not found")
	}
	t.Execute(w, cache.Config)
}

func handleReport(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		token := r.Header.Get("Access-Token")
		var c *cache.UserClient
		c = cache.GetUserClientByToken(token)
		if c == nil {
			c = cache.NewUserClient(r.Header.Get("Referer"), r.Header.Get("Origin"))
			go c.Watch()
		} else {
			c.Refresh()
		}
		w.Header().Set("Access-Token", c.Token)
		// fmt.Println(len(cache.Clients))
	}
	origin := "*"
	if o := r.Header.Get("Origin"); o != "" {
		for _, v := range cache.Config.WhiteList {
			if v == o {
				origin = o
			}
		}
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Expose-Headers", "Access-Token")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Connection, User-Agent, Cookie, Access-Token")
	w.Header().Set("Cookie", "hello cookie")
	w.WriteHeader(http.StatusOK)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
