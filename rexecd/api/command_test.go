package api

import (
	"net/url"
	"strings"
	"testing"
)

func TestBuildCommandQuery(t *testing.T) {
	v := make(url.Values)
	v["username"] = []string{"test_user"}
	v["host_id"] = []string{string("1")}
	got, args, err := buildCommandQuery(v)
	if err != nil {
		t.Fatal(err)
	}
	expected := `
SELECT id, cmd, host_id, timestamp, exit_code
WHERE username = ?
AND host_id = ?
;
  `
	if strings.TrimSpace(expected) != strings.TrimSpace(got) {
		t.Errorf("expected %s, got %s", expected, got)
	}
	username, ok := args[0].(string)
	if !ok {
		t.Error("expected username to be a string")
	}
	if username != "test_user" {
		t.Errorf("expected test_user, got %s", username)
	}

	hostID, ok := args[1].(uint64)
	if !ok {
		t.Error("expected host_id to be a uint64")
	}
	if hostID != 1 {
		t.Errorf("expected test_host, got %d", hostID)
	}
}
