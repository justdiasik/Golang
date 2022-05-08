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

type ShirtResource struct {
	store store.Store
	cache *lru.TwoQueueCache
	broker message_broker.MessageBroker
}

func NewShirtResource(store store.Store, broker message_broker.MessageBroker, cache *lru.TwoQueueCache) *ShirtResource {
	return &ShirtResource{
		store: store,
		broker: broker,
		cache: cache,
	}
}

func (sr *ShirtResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", sr.CreateShirt)
	r.Get("/", sr.AllShirts)
	r.Get("/{id}", sr.ByID)
	r.Put("/", sr.UpdateShirt)
	r.Delete("/{id}", sr.DeleteShirt)

	return r
}

func (sr *ShirtResource) CreateShirt(w http.ResponseWriter, r *http.Request) {
	shirt := new(models.Shirt)
	if err := json.NewDecoder(r.Body).Decode(shirt); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	fmt.Println("Status OK")

	if err := sr.store.Shirts().Create(r.Context(), shirt); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	// Правильно пройтись по всем буквам и всем словам
	sr.broker.Cache().Purge()  // в рамках учебного проекта полностью чистим кэш после создания новой категории

	w.WriteHeader(http.StatusCreated)
}


func (sr *ShirtResource) AllShirts(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.ShirtsFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		shirtsFromCache, ok := sr.cache.Get(searchQuery)
		if ok {
			render.JSON(w, r, shirtsFromCache)
			return
		}

		filter.Query = &searchQuery
	}

	shirts, err := sr.store.Shirts().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if searchQuery != "" {
		sr.cache.Add(searchQuery, shirts)
	}

	render.JSON(w, r, shirts)
}

func (sr *ShirtResource) ByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	shirtFromCache, ok := sr.cache.Get(id)
	if ok {
		render.JSON(w, r, shirtFromCache)
		return
	}

	shirt, err := sr.store.Shirts().ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	sr.cache.Add(id, shirt)
	render.JSON(w, r, shirt)
}

func (sr *ShirtResource) UpdateShirt(w http.ResponseWriter, r *http.Request) {
	shirt := new(models.Shirt)
	if err := json.NewDecoder(r.Body).Decode(shirt); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := sr.store.Shirts().Update(r.Context(), shirt); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	sr.broker.Cache().Remove(shirt.ID)
}

func (sr *ShirtResource) DeleteShirt(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := sr.store.Shirts().Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	sr.broker.Cache().Remove(id)
}
