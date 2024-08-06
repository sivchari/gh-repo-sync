package ghreposync

import (
	"iter"
	"log"
	"os"
	"path"

	"github.com/goccy/go-yaml"
)

type Config struct {
	WorkDir      string   `yaml:"work_dir"`
	Repositories []string `yaml:"repositories"`
}

func Unmarshal(data []byte) (*Config, error) {
	c := &Config{}
	if err := yaml.Unmarshal(data, c); err != nil {
		return nil, err
	}
	return c, nil
}

func Filter(work_dir string, repos []string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, repo := range repos {
			p := path.Join(work_dir, "/", repo)
			stat, err := os.Stat(p)
			if err != nil || !stat.IsDir() {
				log.Println("Error reading file:", err)
				continue
			}
			if !yield(repo) {
				break
			}
		}
	}
}
