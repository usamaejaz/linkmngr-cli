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
	pages := findSubcommand("page", root.Commands())
	if pages == nil {
		t.Fatal("pages command not found")
	}
	if got := findSubcommand("create", pages.Commands()); got != nil {
		t.Fatal("pages create should not exist")
	}
}

func TestLinksCreateFlags(t *testing.T) {
	root := NewRootCmd()
	links := findSubcommand("link", root.Commands())
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

func TestCommandNamingAndAliases(t *testing.T) {
	root := NewRootCmd()

	auth := findSubcommand("auth", root.Commands())
	if auth == nil {
		t.Fatal("auth command not found")
	}
	login := findSubcommand("login", auth.Commands())
	if login == nil || !login.HasAlias("set-token") {
		t.Fatal("auth login command should include set-token alias")
	}
	status := findSubcommand("status", auth.Commands())
	if status == nil || !status.HasAlias("whoami") {
		t.Fatal("auth status command should include whoami alias")
	}
	logout := findSubcommand("logout", auth.Commands())
	if logout == nil || !logout.HasAlias("revoke") {
		t.Fatal("auth logout command should include revoke alias")
	}

	links := findSubcommand("link", root.Commands())
	if links == nil {
		t.Fatal("links command not found")
	}
	if !links.HasAlias("links") {
		t.Fatal("link command should include links alias")
	}
	linkGet := findSubcommand("get", links.Commands())
	if linkGet == nil || !linkGet.HasAlias("view") {
		t.Fatal("links get command should include view alias")
	}
	linkList := findSubcommand("list", links.Commands())
	if linkList == nil || !linkList.HasAlias("ls") {
		t.Fatal("links list command should include ls alias")
	}

	brands := findSubcommand("brand", root.Commands())
	if brands == nil {
		t.Fatal("brands command not found")
	}
	if !brands.HasAlias("brands") {
		t.Fatal("brand command should include brands alias")
	}
	brandGet := findSubcommand("get", brands.Commands())
	if brandGet == nil || !brandGet.HasAlias("view") {
		t.Fatal("brands get command should include view alias")
	}
	domainCheck := findSubcommand("domain-check", brands.Commands())
	if domainCheck == nil || !domainCheck.HasAlias("check-domain") {
		t.Fatal("brands domain-check should include check-domain alias")
	}

	pages := findSubcommand("page", root.Commands())
	if pages == nil {
		t.Fatal("pages command not found")
	}
	if !pages.HasAlias("pages") {
		t.Fatal("page command should include pages alias")
	}
	pageGet := findSubcommand("get", pages.Commands())
	if pageGet == nil || !pageGet.HasAlias("view") {
		t.Fatal("pages get command should include view alias")
	}

	domains := findSubcommand("domain", root.Commands())
	if domains == nil {
		t.Fatal("domain command not found")
	}
	if !domains.HasAlias("domains") {
		t.Fatal("domain command should include domains alias")
	}
	if findSubcommand("list", domains.Commands()) == nil {
		t.Fatal("domains list command not found")
	}
}
