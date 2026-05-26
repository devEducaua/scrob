package db

import (
	"database/sql"
	_ "modernc.org/sqlite"
)
func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./scrob.db");
	if err != nil {
		return nil, err;
	}
	return db, nil;
}

func InitTables(db *sql.DB) error {
	// TODO: reconsider "name" been a unique
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS artists (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL
	);
	`);
	if err != nil {
		return err;
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS albums (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		year INTEGER NOT NULL,
		UNIQUE (title, year)
	);
	`);
	if err != nil {
		return err;
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS tracks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		file TEXT NOT NULL,
		album_id INTEGER NOT NULL REFERENCES albums(id),
		artist_id INTEGER NOT NULL REFERENCES artists(id)
	);
	`);
	if err != nil {
		return err;
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS albums_artists (
		album_id INTEGER NOT NULL REFERENCES albums(id),
		artist_id INTEGER NOT NULL REFERENCES artists(id),
		PRIMARY KEY (album_id, artist_id)
	);
	`);
	if err != nil {
		return err;
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS scrobbles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		time DATETIME DEFAULT CURRENT_TIMESTAMP,
		album_id INTEGER NOT NULL REFERENCES albums(id),
		artist_id INTEGER NOT NULL REFERENCES artists(id),
		track_id INTEGER NOT NULL REFERENCES tracks(id)
	);
	`);
	if err != nil {
		return err;
	}

	return nil;
}

func SaveAlbumArtist(db *sql.DB, albumTitle string, albumYear int, artist string ) (int64, int64, error) {
	albumId, err := SaveAlbum(db, albumTitle, albumYear);
	if err != nil {
		return -1, -1, err;
	}

	artistId, err := SaveArtist(db, artist);
	if err != nil {
		return -1, -1, err;
	}

	_, err = db.Exec(`INSERT OR IGNORE INTO album_artists (album_id, artist_id) VALUES (?, ?)`, albumId, artistId);
	if err != nil {
		return -1, -1, err;
	}
	return albumId, artistId, nil;
}

func SaveAlbum(db *sql.DB, title string, year int) (int64, error) {
	_, err := db.Exec(`INSERT OR IGNORE INTO albums (title, year) VALUES (?, ?)`, title, year);
	if err != nil {
		return -1, err;
	}
	var id int64;
	row := db.QueryRow(`SELECT id FROM albums WHERE title=? AND year=?`, title, year);
	if err := row.Scan(&id); err != nil {
		return -1, err;
	}
	return id, nil;
}

func SaveArtist(db *sql.DB, name string) (int64, error) {
	_, err := db.Exec(`INSERT OR IGNORE INTO artists (name) VALUES (?)`, name);
	if err != nil {
		return -1, err;
	}
	var id int64;
	row := db.QueryRow(`SELECT id FROM artists WHERE name=?`, name);
	if err := row.Scan(&id); err != nil {
		return -1, err;
	}
	return id, nil;
}

func SaveTrack(db *sql.DB, title string, file string, albumId, artistId int64) (int64, error) {
	_, err := db.Exec(`INSERT OR IGNORE INTO tracks (title, file, album_id, artist_id) VALUES (?, ?, ?, ?)`, title, file, albumId, artistId);
	if err != nil {
		return -1, err;
	}
	var id int64;
	row := db.QueryRow(`SELECT id FROM tracks WHERE title=? AND file=? AND album_id? AND artist_id=?`, title, file, albumId, artistId);
	if err := row.Scan(&id); err != nil {
		return -1, err;
	}

	return id, nil;
}

func SaveScrobble(db *sql.DB, albumId, artistId, trackId int64) error {
	_, err := db.Exec(`INSERT INTO scrobbles (album_id, artist_id, track_id) VALUES (?, ?, ?)`, albumId, artistId, trackId);
	if err != nil {
		return err;
	}
	return nil;
}

