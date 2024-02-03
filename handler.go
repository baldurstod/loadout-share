package main

import (
	"encoding/json"
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
			jsonError(w, r, errors.New("Invalid id"))
			return
		}

		loadout, err = getOldLoadout(id)
		if err != nil {
			jsonError(w, r, errors.New("Loadout not found"))
			return
		}

	} else if len(id) == 24 {
		// Mongo format

		matched, err := regexp.Match(`^[a-f0-9]+$`, []byte(id))
		if err != nil || !matched {
			jsonError(w, r, errors.New("Invalid id"))
			return
		}

		loadout, err = getLoadout(id)
		if err != nil {
			jsonError(w, r, errors.New("Loadout not found"))
			return
		}

	} else {
		jsonError(w, r, errors.New("Invalid id"))
		return
	}

	jsonSuccess(w, r, loadout)
}

func (h *Handler) servePost(w http.ResponseWriter, r *http.Request) {
	var body = make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		jsonError(w, r, errors.New("Bad request"))
		return
	}

	loadoutID, err := addLoadout(body)
	if err != nil {
		log.Println(err)
		jsonError(w, r, errors.New("Error while creating loadout"))
		return
	}

	jsonSuccess(w, r, loadoutID)
}
