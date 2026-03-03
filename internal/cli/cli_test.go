package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func findSubcommand(cmdName string, parentCommands []*cobra.Command) *cobra.Command {
	for _, c := range parentCommands {
		if c.Name() == cmdName {
			return c
		}
	}
	return nil
}

func TestPagesCreateCommandRemoved(t *testing.T) {
	root := NewRootCmd()
	pages := findSubcommand("pages", root.Commands())
	if pages == nil {
		t.Fatal("pages command not found")
	}
	if got := findSubcommand("create", pages.Commands()); got != nil {
		t.Fatal("pages create should not exist")
	}
}

func TestLinksCreateFlags(t *testing.T) {
	root := NewRootCmd()
	links := findSubcommand("links", root.Commands())
	if links == nil {
		t.Fatal("links command not found")
	}
	create := findSubcommand("create", links.Commands())
	if create == nil {
		t.Fatal("links create command not found")
	}

	for _, flagName := range []string{"domain", "slug", "brand-id"} {
		if create.Flags().Lookup(flagName) == nil {
			t.Fatalf("expected %q flag on links create", flagName)
		}
	}
}

func TestBuildRequestBody(t *testing.T) {
	dir := t.TempDir()
	dataFile := filepath.Join(dir, "body.json")
	if err := os.WriteFile(dataFile, []byte(`{"a":"b"}`), 0o600); err != nil {
		t.Fatalf("write data file: %v", err)
	}

	body, err := buildRequestBody("", dataFile, nil)
	if err != nil {
		t.Fatalf("buildRequestBody from file returned error: %v", err)
	}
	if body["a"] != "b" {
		t.Fatalf("expected field from file, got: %#v", body)
	}

	body, err = buildRequestBody("", "", []string{"title=My Link", "slug=my-link"})
	if err != nil {
		t.Fatalf("buildRequestBody from --set returned error: %v", err)
	}
	if body["title"] != "My Link" || body["slug"] != "my-link" {
		t.Fatalf("unexpected body from --set: %#v", body)
	}
}
