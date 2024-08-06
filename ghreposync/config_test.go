package ghreposync_test

import (
	_ "embed"
	"os"
	"reflect"
	"testing"

	"github.com/sivchari/gh-repo-sync/ghreposync"
)

//go:embed testdata/gh-repo-sync.yaml
var testdata []byte

func TestUnmarshal(t *testing.T) {
	want := &ghreposync.Config{
		WorkDir: "~/workspace",
		Repositories: []string{
			"aaaa/xxxx",
			"bbbb/yyyy",
			"cccc/zzzz",
		},
	}
	got, err := ghreposync.Unmarshal(testdata)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want: %v, got: %v", want, got)
	}
}

func TestFilter(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repos := []string{
		"testdata",
		"undefined",
		"testdata/gh-repo-sync.yaml",
	}
	want := []string{"testdata"}
	got := []string{}
	for repo := range ghreposync.Filter(wd, repos) {
		got = append(got, repo)
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want: %v, got: %v", want, got)
	}
}
