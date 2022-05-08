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

type LaptopResource struct {
	store store.Store
	cache *lru.TwoQueueCache
	broker message_broker.MessageBroker
}

func NewLaptopResource(store store.Store, broker message_broker.MessageBroker, cache *lru.TwoQueueCache) *LaptopResource {
	return &LaptopResource{
		store: store,
		broker: broker,
		cache: cache,
	}
}

func (lr *LaptopResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", lr.CreateLaptop)
	r.Get("/", lr.AllLaptops)
	r.Get("/{id}", lr.ByID)
	r.Put("/", lr.UpdateLaptop)
	r.Delete("/{id}", lr.DeleteLaptop)

	return r
}

func (lr *LaptopResource) CreateLaptop(w http.ResponseWriter, r *http.Request) {
	laptop := new(models.Laptop)
	if err := json.NewDecoder(r.Body).Decode(laptop); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	fmt.Println("Status OK")

	if err := lr.store.Laptops().Create(r.Context(), laptop); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	// Правильно пройтись по всем буквам и всем словам
	lr.broker.Cache().Purge()  // в рамках учебного проекта полностью чистим кэш после создания новой категории

	w.WriteHeader(http.StatusCreated)
}


func (lr *LaptopResource) AllLaptops(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.LaptopsFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		laptopsFromCache, ok := lr.cache.Get(searchQuery)
		if ok {
			render.JSON(w, r, laptopsFromCache)
			return
		}

		filter.Query = &searchQuery
	}

	laptops, err := lr.store.Laptops().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if searchQuery != "" {
		lr.cache.Add(searchQuery, laptops)
	}

	render.JSON(w, r, laptops)
}

func (lr *LaptopResource) ByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	laptopFromCache, ok := lr.cache.Get(id)
	if ok {
		render.JSON(w, r, laptopFromCache)
		return
	}

	laptop, err := lr.store.Laptops().ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	lr.cache.Add(id, laptop)
	render.JSON(w, r, laptop)
}

func (lr *LaptopResource) UpdateLaptop(w http.ResponseWriter, r *http.Request) {
	laptop := new(models.Laptop)
	if err := json.NewDecoder(r.Body).Decode(laptop); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := lr.store.Laptops().Update(r.Context(), laptop); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	lr.broker.Cache().Remove(laptop.ID)
}

func (lr *LaptopResource) DeleteLaptop(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := lr.store.Laptops().Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	lr.broker.Cache().Remove(id)
}