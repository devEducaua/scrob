package scrobble

import (
	"scrob/internal/db"
	"scrob/internal/mpd"
)

func ScrobbleCurrentSong() (error) {
	s, err := mpd.GetCurrentSong();
	if err != nil {
		return err;
	}
	
	connDB, err := db.Connect();
	if err != nil {
		return err;
	}
	defer connDB.Close();

	albumId, artistId, err := db.SaveAlbumArtist(connDB, s.Album, s.Date, s.Artist);
	if err != nil {
		return err;
	}

	trackId, err := db.SaveTrack(connDB, s.Title, s.File, albumId, artistId);
	if err != nil {
		return err;
	}

	err = db.SaveScrobble(connDB, albumId, artistId, trackId);
	if err != nil {
		return err;
	}

	return nil;
}

func WatchPlayer(msg chan<- string, errs chan<- error) {
	for {
		result, err := mpd.Request("idle player");
		if err != nil {
			errs <- err;
			continue;
		}
		msg <- result;
	}
}
