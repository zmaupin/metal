package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// CommandResponse is the response from a command
type CommandResponse struct {
	Command []CommandJSON
}

// CommandJSON is an item in a CommandResponse
type CommandJSON struct {
	ID        uint64 `json:"id"`
	Cmd       string `json:"cmd"`
	HostID    uint64 `json:"host_id"`
	Timestamp uint64 `json:"timestamp"`
	ExitCode  int    `json:"exit_code"`
}

func command(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	response := CommandResponse{Command: []CommandJSON{}}
	query, args, err := buildCommandQuery(r)
	if err != nil {
		log.Error(err)
		w.WriteHeader(500)
		return
	}
	rows, err := db.Query(query, args...)
	defer rows.Close()
	if err != nil {
		w.WriteHeader(500)
		log.Error(err)
		return
	}
	for rows.Next() {
		var id uint64
		var cmd string
		var hostID uint64
		var timestamp uint64
		var exitCode int

		if err = rows.Scan(id, cmd, hostID, timestamp, exitCode); err != nil {
			log.Error(err)
			w.WriteHeader(500)
			return
		}
		response.Command = append(response.Command, CommandJSON{
			ID:        id,
			Cmd:       cmd,
			HostID:    hostID,
			Timestamp: timestamp,
			ExitCode:  exitCode,
		})
	}
	b, err := json.Marshal(&response)
	if err != nil {
		w.WriteHeader(500)
		log.Error(err)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		w.WriteHeader(500)
		log.Error(err)
		return
	}
}

func buildCommandQuery(r *http.Request) (query string, args []interface{}, err error) {
	args = []interface{}{}
	query = `
  SELECT id, cmd, host_id, timestamp, exit_code
  `
	if idStr, ok := r.Form["id"]; ok {
		id, err := strconv.Atoi(idStr[0])
		if err != nil {
			return query, args, err
		}
		query += `
    WHERE id = ?
    `
		args = append(args, id)
	}

	if cmd, ok := r.Form["cmd"]; ok {
		query += `
    WHERE cmd = ?
    `
		args = append(args, cmd)
	}

	if username, ok := r.Form["username"]; ok {
		query += `
    WHERE username = ?
    `
		args = append(args, username)
	}

	if hostID, ok := r.Form["host_id"]; ok {
		query += `
    WHERE host_id = ?
    `
		args = append(args, hostID)
	}

	if exitCode, ok := r.Form["exit_code"]; ok {
		query += `
    WHERE exit_code = ?
    `
		args = append(args, exitCode)
	}
	query += ";"
	return query, args, nil
}
