package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"scrob/internal/api"
	"scrob/internal/db"
	"scrob/internal/mpd"
	"scrob/internal/scrobble"
)

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

	go serve(":7087");

	msgs := make(chan string);
	errs := make(chan error);

	go watchPlayer(msgs, errs);

	for {
		select {
		case msg := <-msgs:
			if msg == "changed: player\n" {
				if err := scrobble.ScrobbleCurrentSong(); err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: %v\n", err);
				}
			}
		case err := <-errs:
			fatal(err);
		}
	}
}

func serve(port string) {
	listener, err := net.Listen("tcp", port);
	if err != nil {
		fatal(err);
	}

	for {
		conn, err := listener.Accept();
		if err != nil {
			fmt.Fprintln(os.Stderr, err);
		}
		go handleConnection(conn);
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close();

	reader := bufio.NewReader(conn);
	req, err := reader.ReadString('\n');
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err);
	}

	resp := api.ParseRequest(req);
	bt, err := json.MarshalIndent(resp, "", "    ");
	if err != nil {
		fatal(err);
	}
	conn.Write(bt);
}

func watchPlayer(msg chan<- string, errs chan<- error) {
	for {
		result, err := mpd.Request("idle player");
		if err != nil {
			errs <- err;
			continue;
		}
		msg <- result;
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "ERROR: %v\n", err);
	os.Exit(1);
}
