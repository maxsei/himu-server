package main

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

//go:embed schema.sql
var DbSchema string

func init() {
	log = logrus.New()
	var formatter logrus.TextFormatter
	formatter.ForceColors = true
	formatter.FullTimestamp = true
	log.Formatter = &formatter
}

func main() {
	db, err := sql.Open("sqlite3", "file:./data/sensor_data.sqlite3?cache=shared")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	{
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if _, err := db.ExecContext(ctx, DbSchema); err != nil {
			log.Panic(err)
		}
	}

	port := "2055"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	ln, err := net.Listen("tcp", net.JoinHostPort("", port))
	if err != nil {
		log.Panic(err)
	}
	log.Infof("listening to tcp connection on port %s", port)
	serverStart := time.Now().UnixMilli()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Error(err)
			return
		}
		log.Infof("got connection %s", conn.LocalAddr())
		go handle(conn, db, serverStart)
	}
}

func handle(conn net.Conn, db *sql.DB, serverStart int64) {
	defer conn.Close()
	defer log.Warn("closed connection")
	dec := json.NewDecoder(conn)
outer:
	for {
		if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
			log.Error(err)
			continue
		}
		// if _, err := io.Copy(os.Stdout, buf); err != nil {
		// 	log.Error(err)
		// 	return
		// }
		seq := []string{
			"{",
			"os",
			"hyperimu",
		}
		for _, s := range seq {
			t, err := dec.Token()
			if err != nil {
				log.Error(err)
				return
			}
			if s != fmt.Sprintf("%v", t) {
				continue outer
			}
		}
		for dec.More() {
			sensor, err := dec.Token()
			if err != nil {
				log.Error(err)
				return
			}
			var values []float32
			if err := dec.Decode(&values); err != nil {
				log.Error(err)
				return
			}
			tsms := time.Now().UnixMilli()
			log.Infof("%s: %v", sensor, values)
			InsertRecord(db, &Record{
				serverStart: serverStart,
				sensor: fmt.Sprint(sensor),
				tsms:   tsms,
				values: values,
			})
		}
	}
}

type Record struct {
	serverStart   int64
	sensor string
	tsms   int64
	values []float32
}

func InsertRecord(db_ *sql.DB, r *Record) error {
	db, err := db_.Begin()
	if err != nil {
		return err
	}
	defer db.Rollback()
	for coord, value := range r.values {
		if _, err := db.Exec(`
				INSERT INTO sensor_data (
					server_start_ts_ms, ts_ms, sensor, value, coord
				) VALUES (
					?, ?, ?, ?, ?
				);
			`, r.serverStart, r.tsms, r.sensor, value, coord); err != nil {
			return err
		}
	}
	return db.Commit()
}
