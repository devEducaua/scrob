package mpd

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Song struct {
	File string
	LastModified string
	Added string
	Format string
	Artist string
	Title string
	Album string
	Date int
	Time int
	Duration float64
	Position int
	Id int
}

func initializeMpdConnection() (*net.TCPConn, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:6600");
	if err != nil {
		return nil, err;
	}

	conn, err := net.DialTCP("tcp", nil, addr);
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mpd: %v", err);
	}

	reader := bufio.NewReader(conn);
	line, err := reader.ReadString('\n');
	if err != nil {
		return nil, err;
	}

	if !strings.HasPrefix(line, "OK MPD") {
		return nil, fmt.Errorf("failed to initialize mpd connection");
	}

	return conn, nil;
}

func Request(request string) (string, error) {
	conn, err := initializeMpdConnection();
	if err != nil {
		return "", err;
	}
	if conn != nil {
		defer conn.Close();
	}

	var line string;
	var reader = bufio.NewReader(conn);
	var sb strings.Builder;

	fmt.Fprintf(conn, "%v\n", request);
	loop:
	for {
		if line, err = reader.ReadString('\n'); err != nil {
			return "",  err;
		}

		switch {
		case line == "OK\n":
			break loop;
		case strings.HasPrefix(line, "ACK "):
			err = fmt.Errorf("request failed: %v", line);
			break loop;
		default:
			sb.WriteString(line);
		}
	}

	return sb.String(), err;
}

func GetCurrentSong() (Song, error) {
	result, err := Request("currentsong");
	if err != nil {
		return Song{}, err;
	}

	s, err := ParseSongResponse(result);
	if err != nil {
		return Song{}, err;
	}
	return s, nil;
}

func ParseSongResponse(resp string) (Song, error) {
	var s Song;
	lines := strings.Split(resp, "\n");
	for i := range lines {
		line := strings.TrimSpace(lines[i]);

		switch {
		case strings.HasPrefix(line, "file: "):
			s.File = line[5:];
		case strings.HasPrefix(line, "Last-Modified: "):
			s.LastModified = line[15:];
		case strings.HasPrefix(line, "Added: "):
			s.Added = line[7:];
		case strings.HasPrefix(line, "Artist: "):
			s.Artist = line[8:];
		case strings.HasPrefix(line, "Title: "):
			s.Title = line[7:];
		case strings.HasPrefix(line, "Album: "):
			s.Album = line[7:];
		case strings.HasPrefix(line, "Date: "):
			str := line[6:];
			converted, err := strconv.Atoi(str);
			if err != nil {
				return Song{}, err;
			}
			s.Date = converted;
		case strings.HasPrefix(line, "Time: "):
			str := line[6:];
			converted, err := strconv.Atoi(str);
			if err != nil {
				return Song{}, err;
			}
			s.Time = converted;
		case strings.HasPrefix(line, "Duration: "):
			str := line[10:];
			converted, err := strconv.ParseFloat(str, 64);
			if err != nil {
				return Song{}, err;
			}
			s.Duration = converted;
		case strings.HasPrefix(line, "Position: "):
			str := line[10:];
			converted, err := strconv.Atoi(str);
			if err != nil {
				return Song{}, err;
			}
			s.Position = converted;
		case strings.HasPrefix(line, "Id: "):
			str := line[4:];
			converted, err := strconv.Atoi(str);
			if err != nil {
				return Song{}, err;
			}
			s.Id = converted;
		}
	}
	return s, nil;
}
