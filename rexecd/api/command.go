package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	query, args, err := buildCommandQuery(r.Form)
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

func buildCommandQuery(v url.Values) (query string, args []interface{}, err error) {
	args = []interface{}{}
	query = `
SELECT id, cmd, host_id, timestamp, exit_code
`

	condition := func(key string) string {
		if len(args) == 0 {
			return fmt.Sprintf("WHERE %s = ?\n", key)
		}
		return fmt.Sprintf("AND %s = ?\n", key)
	}

	if idStr, ok := v["id"]; ok {
		id, err := strconv.Atoi(idStr[0])
		if err != nil {
			return query, args, err
		}
		query += condition("id")
		args = append(args, uint64(id))
	}

	if cmd, ok := v["cmd"]; ok {
		query += condition("cmd")
		args = append(args, cmd[0])
	}

	if username, ok := v["username"]; ok {
		query += condition("username")
		args = append(args, username[0])
	}

	if hostID, ok := v["host_id"]; ok {
		id, err := strconv.Atoi(hostID[0])
		if err != nil {
			return query, args, err
		}
		query += condition("host_id")
		args = append(args, uint64(id))
	}

	if exitCode, ok := v["exit_code"]; ok {
		query += condition("exit_code")
		args = append(args, exitCode[0])
	}
	query += ";"
	return query, args, nil
}
