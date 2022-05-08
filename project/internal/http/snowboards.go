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

type SnowboardResource struct {
	store store.Store
	broker message_broker.MessageBroker
	cache *lru.TwoQueueCache
}

func NewSnowboardResource(store store.Store, broker message_broker.MessageBroker, cache *lru.TwoQueueCache) *SnowboardResource {
	return &SnowboardResource{
		store: store,
		broker: broker,
		cache: cache,
	}
}

func (sr *SnowboardResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", sr.CreateSnowboard)
	r.Get("/", sr.AllSnowboards)
	r.Get("/{id}", sr.ByID)
	r.Put("/", sr.UpdateSnowboard)
	r.Delete("/{id}", sr.DeleteSnowboard)

	return r
}

func (sr *SnowboardResource) CreateSnowboard(w http.ResponseWriter, r *http.Request) {
	snowboard := new(models.Snowboard)
	if err := json.NewDecoder(r.Body).Decode(snowboard); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	fmt.Println("Status OK")

	if err := sr.store.Snowboards().Create(r.Context(), snowboard); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	// Правильно пройтись по всем буквам и всем словам
	sr.broker.Cache().Purge()  // в рамках учебного проекта полностью чистим кэш после создания новой категории

	w.WriteHeader(http.StatusCreated)
}


func (sr *SnowboardResource) AllSnowboards(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.SnowboardsFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		snowboardsFromCache, ok := sr.cache.Get(searchQuery)
		if ok {
			render.JSON(w, r, snowboardsFromCache)
			return
		}

		filter.Query = &searchQuery
	}

	snowboards, err := sr.store.Snowboards().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if searchQuery != "" {
		sr.cache.Add(searchQuery, snowboards)
	}

	render.JSON(w, r, snowboards)
}

func (sr *SnowboardResource) ByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	snowboardFromCache, ok := sr.cache.Get(id)
	if ok {
		render.JSON(w, r, snowboardFromCache)
		return
	}

	snowboard, err := sr.store.Snowboards().ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	sr.cache.Add(id, snowboard)
	render.JSON(w, r, snowboard)
}

func (sr *SnowboardResource) UpdateSnowboard(w http.ResponseWriter, r *http.Request) {
	snowboard := new(models.Snowboard)
	if err := json.NewDecoder(r.Body).Decode(snowboard); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := sr.store.Snowboards().Update(r.Context(), snowboard); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	sr.broker.Cache().Remove(snowboard.ID)
}

func (sr *SnowboardResource) DeleteSnowboard(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := sr.store.Snowboards().Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	sr.broker.Cache().Remove(id)
}
