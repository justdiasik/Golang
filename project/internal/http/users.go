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
)


type UserResource struct {
	store store.Store
	broker message_broker.MessageBroker
	cache *lru.TwoQueueCache
}

func NewUserResource(store store.Store, broker message_broker.MessageBroker, cache *lru.TwoQueueCache) *UserResource {
	return &UserResource{
		store: store,
		broker: broker,
		cache: cache,
	}
}

func (ur *UserResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", ur.CreateUser)
	r.Get("/", ur.AllUsers)

	return r
}

func (ur *UserResource) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := new(models.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}



	//if !user.IsEmailValid() {
	//	fmt.Fprintln(w, "Invalid email error!")
	//	return
	//}
	//if !user.IsPasswordValid() {
	//	fmt.Fprintln(w, "Invalid password error!")
	//	return
	//}
	if !user.IsEmailValid() || !user.IsPasswordValid() {
		fmt.Fprintln(w, "Invalid email or password!")
		return
	}
	fmt.Fprintln(w, "Success!")

	if err := ur.store.Users().Create(r.Context(), user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	// Правильно пройтись по всем буквам и всем словам
	ur.broker.Cache().Purge()  // в рамках учебного проекта полностью чистим кэш после создания новой категории

	w.WriteHeader(http.StatusCreated)
}


func (ur *UserResource) AllUsers(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.UsersFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		usersFromCache, ok := ur.cache.Get(searchQuery)
		if ok {
			render.JSON(w, r, usersFromCache)
			return
		}

		filter.Query = &searchQuery
	}

	users, err := ur.store.Users().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if searchQuery != "" {
		ur.cache.Add(searchQuery, users)
	}

	render.JSON(w, r, users)
}




