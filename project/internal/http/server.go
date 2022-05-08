package http

import (
	"context"
	"github.com/go-chi/chi"
	lru "github.com/hashicorp/golang-lru"
	"lectures-6/internal/message_broker"
	"lectures-6/internal/store"
	"log"
	"net/http"
	"time"
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}
	store       store.Store
	cache 		*lru.TwoQueueCache
	broker      message_broker.MessageBroker

	Address string
}

func NewServer(ctx context.Context, opts ...ServerOption) *Server {
	srv := &Server{
		ctx:         ctx,
		idleConnsCh: make(chan struct{}),
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()

	laptopsResource := NewLaptopResource(s.store, s.broker, s.cache)
	r.Mount("/electronics/laptops", laptopsResource.Routes())

	snowboardsResource := NewSnowboardResource(s.store, s.broker, s.cache)
	r.Mount("/sport-hobby/snowboards", snowboardsResource.Routes())

	shirtsResource := NewShirtResource(s.store, s.broker, s.cache)
	r.Mount("/fashion-style/shirts", shirtsResource.Routes())

	toysResource := NewToyResource(s.store, s.broker, s.cache)
	r.Mount("/for-kids/toys", toysResource.Routes())

	usersResource := NewUserResource(s.store, s.broker, s.cache)
	r.Mount("/registration", usersResource.Routes())

	return r
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:         s.Address,
		Handler:      s.basicHandler(),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 30,
	}
	go s.ListenCtxForGT(srv)

	log.Println("[HTTP] Server running on", s.Address)
	return srv.ListenAndServe()
}

func (s *Server) ListenCtxForGT(srv *http.Server) {
	<-s.ctx.Done() // блокируемся, пока контекст приложения не отменен

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("[HTTP] Got err while shutting down^ %v", err)
	}

	log.Println("[HTTP] Proccessed all idle connections")
	close(s.idleConnsCh)
}

func (s *Server) WaitForGracefulTermination() {
	// блок до записи или закрытия канала
	<-s.idleConnsCh
}
