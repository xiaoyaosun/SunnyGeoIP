package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	signal.Notify(sigchan, syscall.SIGTERM)

	server := NewHttpServer()
	server.listen(":8087")

	server.bind(NewGeolocationServer())

	<-sigchan
	server.close()

}
