package main

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"runtime/debug"
)

type Handler struct {
	Handler http.Handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, string(debug.Stack()))
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	switch r.Method {
	case "GET":
		h.serveGet(w, r)
	case "POST":
		h.servePost(w, r)
	}
}

func (h *Handler) serveGet(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[1:]
	var loadout map[string]interface{}

	if len(id) == 16 {
		// NEDB format

		matched, err := regexp.Match(`^[a-zA-Z0-9]+$`, []byte(id))
		if err != nil || !matched {
			jsonError(w, errors.New("Invalid id"))
			return
		}

		loadout, err = getOldLoadout(id)
		if err != nil {
			jsonError(w, errors.New("Loadout not found"))
			return
		}

	} else if len(id) == 24 {
		// Mongo format

		matched, err := regexp.Match(`^[a-f0-9]+$`, []byte(id))
		if err != nil || !matched {
			jsonError(w, errors.New("Invalid id"))
			return
		}

		loadout, err = getLoadout(id)
		if err != nil {
			jsonError(w, errors.New("Loadout not found"))
			return
		}

	} else {
		jsonError(w, errors.New("Invalid id"))
	}

	jsonSuccess(w, loadout)
}

func (h *Handler) servePost(w http.ResponseWriter, r *http.Request) {
}
