package main

import 	(
	"context"
	lru "github.com/hashicorp/golang-lru"
	"lectures-6/internal/http"
	"lectures-6/internal/message_broker/kafka"
	"lectures-6/internal/store/postgres"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go CatchTermination(cancel)

	dbURL := "postgres://localhost:5431/laptopsbd"
	store := postgres.NewDB()
	if err := store.Connect(dbURL); err != nil {
		panic(err)
	}
	defer store.Close()

	cache, err := lru.New2Q(6)
	if err != nil {
		panic(err)
	}

	brokers := []string{"localhost:29092"}
	broker := kafka.NewBroker(brokers, cache, "peer2")
	if err := broker.Connect(ctx); err != nil {
		panic(err)
	}
	defer broker.Close()


	srv := http.NewServer(
		ctx,
		http.WithAddress(":8075"),
		http.WithStore(store),
		http.WithCache(cache),
		http.WithBroker(broker),
	)
	if err := srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGracefulTermination()
}

func CatchTermination(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Print("[WARN] caught termination signal")
	cancel()
}
