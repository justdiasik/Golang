package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	lru "github.com/hashicorp/golang-lru"
	"lectures-6/internal/message_broker"
	"lectures-6/internal/models"
	"lectures-6/internal/store"
	"net/http"
	"strconv"
)

type ToyResource struct {
	store store.Store
	broker message_broker.MessageBroker
	cache *lru.TwoQueueCache
}

func NewToyResource(store store.Store, broker message_broker.MessageBroker, cache *lru.TwoQueueCache) *ToyResource {
	return &ToyResource{
		store: store,
		broker: broker,
		cache: cache,
	}
}

func (tr *ToyResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", tr.CreateToy)
	r.Get("/", tr.AllToys)
	r.Get("/{id}", tr.ByID)
	r.Put("/", tr.UpdateToy)
	r.Delete("/{id}", tr.DeleteToy)

	return r
}

func (tr *ToyResource) CreateToy(w http.ResponseWriter, r *http.Request) {
	toy := new(models.Toy)
	if err := json.NewDecoder(r.Body).Decode(toy); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	fmt.Println("Status OK")

	if err := tr.store.Toys().Create(r.Context(), toy); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	// Правильно пройтись по всем буквам и всем словам
	tr.broker.Cache().Purge()  // в рамках учебного проекта полностью чистим кэш после создания новой категории

	w.WriteHeader(http.StatusCreated)
}


func (tr *ToyResource) AllToys(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.ToysFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		toysFromCache, ok := tr.cache.Get(searchQuery)
		if ok {
			render.JSON(w, r, toysFromCache)
			return
		}

		filter.Query = &searchQuery
	}

	toys, err := tr.store.Toys().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if searchQuery != "" {
		tr.cache.Add(searchQuery, toys)
	}

	render.JSON(w, r, toys)
}

func (tr *ToyResource) ByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	toyFromCache, ok := tr.cache.Get(id)
	if ok {
		render.JSON(w, r, toyFromCache)
		return
	}

	toy, err := tr.store.Toys().ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	tr.cache.Add(id, toy)
	render.JSON(w, r, toy)
}

func (tr *ToyResource) UpdateToy(w http.ResponseWriter, r *http.Request) {
	toy := new(models.Toy)
	if err := json.NewDecoder(r.Body).Decode(toy); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := tr.store.Toys().Update(r.Context(), toy); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	tr.broker.Cache().Remove(toy.ID)
}

func (tr *ToyResource) DeleteToy(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := tr.store.Toys().Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	tr.broker.Cache().Remove(id)
}
