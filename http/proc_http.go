package http

import (
	cutils "github.com/open-falcon/common/utils"
	"github.com/open-falcon/graph/proc"
	"github.com/open-falcon/graph/store"
	"net/http"
	"strconv"
	"strings"
)

func configProcRoutes() {
	// statistics, TO BE DISCARDED
	http.HandleFunc("/statistics/all", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, proc.GetAll())
	})

	http.HandleFunc("/counter/all", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, proc.GetAll())
	})

	// trace, TO BE DISCARDED
	http.HandleFunc("/trace/", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Path[len("/trace/"):]
		args := strings.Split(urlParam, "/")

		argsLen := len(args)
		endpoint := args[0]
		metric := args[1]
		tags := make(map[string]string)
		if argsLen > 2 {
			tagVals := strings.Split(args[2], ",")
			for _, tag := range tagVals {
				tagPairs := strings.Split(tag, "=")
				if len(tagPairs) == 2 {
					tags[tagPairs[0]] = tagPairs[1]
				}
			}
		}
		proc.RecvDataTrace.SetPK(cutils.Checksum(endpoint, metric, tags))
		RenderDataJson(w, proc.RecvDataTrace.GetAllTraced())
	})

	// filter
	http.HandleFunc("/filter/", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Path[len("/filter/"):]
		args := strings.Split(urlParam, "/")

		argsLen := len(args)
		endpoint := args[0]
		metric := args[1]
		opt := args[2]

		threadholdStr := args[3]
		threadhold, err := strconv.ParseFloat(threadholdStr, 64)
		if err != nil {
			RenderDataJson(w, "bad threadhold")
			return
		}

		tags := make(map[string]string)
		if argsLen > 4 {
			tagVals := strings.Split(args[4], ",")
			for _, tag := range tagVals {
				tagPairs := strings.Split(tag, "=")
				if len(tagPairs) == 2 {
					tags[tagPairs[0]] = tagPairs[1]
				}
			}
		}

		err = proc.RecvDataFilter.SetFilter(cutils.Checksum(endpoint, metric, tags), opt, threadhold)
		if err != nil {
			RenderDataJson(w, err.Error())
			return
		}

		RenderDataJson(w, proc.RecvDataFilter.GetAllFiltered())
	})

	// history
	http.HandleFunc("/history/", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Path[len("/history/"):]
		args := strings.Split(urlParam, "/")

		argsLen := len(args)
		endpoint := args[0]
		metric := args[1]
		tags := make(map[string]string)
		if argsLen > 2 {
			tagVals := strings.Split(args[2], ",")
			for _, tag := range tagVals {
				tagPairs := strings.Split(tag, "=")
				if len(tagPairs) == 2 {
					tags[tagPairs[0]] = tagPairs[1]
				}
			}
		}
		RenderDataJson(w, store.GetAllItems(cutils.Checksum(endpoint, metric, tags)))
	})

	// last
	http.HandleFunc("/last/", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Path[len("/last/"):]
		args := strings.Split(urlParam, "/")

		argsLen := len(args)
		endpoint := args[0]
		metric := args[1]
		tags := make(map[string]string)
		if argsLen > 2 {
			tagVals := strings.Split(args[2], ",")
			for _, tag := range tagVals {
				tagPairs := strings.Split(tag, "=")
				if len(tagPairs) == 2 {
					tags[tagPairs[0]] = tagPairs[1]
				}
			}
		}
		RenderDataJson(w, store.GetLastItem(cutils.Checksum(endpoint, metric, tags)))
	})

}
