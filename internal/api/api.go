package api

import (
	"fmt"
	"scrob/internal/db"
	"strings"
)

/*
FIELD := track/artist/album

rank FIELD [LIMIT]: most scrobbled
list FIELD [LIMIT]: last scrobbles
get FIELD ID [LIMIT]: scrobbles by specific target
*/

type Response struct {
	Ok bool `json:"ok"`
	Error string `json:"error"`
	Data any `json:"data"`
}

func ParseRequest(msg string) (Response) {
	r := Response{
		Ok: true,
	};

	if len(msg) == 0 {
		r.Ok = false;
		r.Error = "request could not be null";
		return r;
	}	

	parts := strings.Split(msg, " ");
	if len(parts) < 2 {
		r.Ok = false;
		r.Error = "invalid number of parameters in the request";
		return r;
	}
	connDB, err := db.Connect();	
	if err != nil {
		r.Ok = false;
		r.Error = fmt.Sprintf("could not connect to the database: %v", err);
		return r;
	}
	defer connDB.Close();

	command := parts[0];
	switch command {
	case "rank":
		subcommand := strings.TrimSpace(parts[1]);
		switch subcommand {
		case "album":
			albums, err := db.GetAlbumsRank(connDB);
			if err != nil {
				r.Ok = false;
				r.Error = err.Error();
				return r;
			}
			r.Data = albums;
		case "track":
			tracks, err := db.GetTracksRank(connDB);
			if err != nil {
				r.Ok = false;
				r.Error = err.Error();
				return r;
			}
			r.Data = tracks;
		}
	case "list":
	case "get":
	default:
		r.Ok = false;
		r.Error = fmt.Sprintf("unknown command: %v", command);
		return r;
	}

	return r;
}
