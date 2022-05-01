package main

import "fmt"

func a(ch chan<- bool) {
	for i := 0; i < 50; i++ {
		fmt.Print("a")
	}
	ch <- true
}

func b(ch chan<- bool) {
	for i := 0; i < 50; i++ {
		fmt.Print("b")
	}
	ch <- true
}

func main() {
	ch := make(chan bool)
	go a(ch)
	go b(ch)
	<-ch
	<-ch
	fmt.Println("end main()")
}
