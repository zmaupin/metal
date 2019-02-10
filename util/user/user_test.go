// +build unit

package user

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestHome(t *testing.T) {
	got, err := Home()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(runtime.GOOS, "darwin") {
		expected := fmt.Sprintf("/Users/%s", os.Getenv("USER"))
		if expected != got {
			t.Errorf("expected %s, got %s", expected, got)
		}
	} else {
		expected := fmt.Sprintf("/home/%s", os.Getenv("USER"))
		if expected != got {
			t.Errorf("expected%s, got %s", expected, got)
		}
	}
}

func TestExpand(t *testing.T) {
	got, err := Expand("~/thing")
	if err != nil {
		t.Fatal(err)
	}
	switch runtime.GOOS {
	case "darwin":
		expected := fmt.Sprintf("/Users/%s/thing", os.Getenv("USER"))
		if expected != got {
			t.Errorf("expected %s, got %s", expected, got)
		}

	case "linux":
		expected := fmt.Sprintf("/home/%s/thing", os.Getenv("USER"))
		if expected != got {
			t.Errorf("expected %s, got %s", expected, got)
		}
	}
}
