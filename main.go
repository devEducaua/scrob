package main

import (
	"fmt"
	"os"
	"scrob/internal/db"
	"scrob/internal/mpd"
	"scrob/internal/scrobble"
)

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "ERROR: %v\n", err);
	os.Exit(1);
}

func main() {

	connDB, err := db.Connect();
	if err != nil {
		fatal(err);
	}
	defer connDB.Close();

	err = db.InitTables(connDB);
	if err != nil {
		fatal(err);
	}

	msgs := make(chan string);
	errs := make(chan error);

	go watchPlayer(msgs, errs);

	for {
		select {
		case msg := <-msgs:
			if msg == "changed: player\n" {
				if err := scrobble.ScrobbleCurrentSong(); err != nil {
					fmt.Fprintln(os.Stderr, err);
				}
			}
		case err := <-errs:
			fatal(err);
		}
	}
}

func watchPlayer(msg chan<- string, errs chan<- error) {
	for {
		result, err := mpd.Request("idle player");
		if err != nil {
			errs <- err;
		}
		msg <- result;
	}
}
