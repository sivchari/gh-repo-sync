package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/sivchari/gh-repo-sync/ghreposync"
	"github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
)

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
}

func main() {
	_main()
}

func _main() {
	var file *string
	var timeout *time.Duration
	file = pflag.StringP("file", "f", "gh-repo-sync.yaml", "config file")
	timeout = pflag.DurationP("timeout", "t", 5*time.Minute, "timeout for each operation")

	pflag.Usage = func() {
		fmt.Printf(`
Usage:
  gh repo-sync [flags]

Flags:
%s
Examples:
  gh repo-sync -f gh-repo-sync.yaml
`, pflag.CommandLine.FlagUsages())
	}
	pflag.Parse()

	if *file == "" {
		pflag.Usage()
		return
	}

	f, err := os.ReadFile(*file)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	cfg, err := ghreposync.Unmarshal(f)
	if err != nil {
		log.Fatal("Error unmarshalling config:", err)
	}

	signalCtx, signalCancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer signalCancel()

	timeoutCtx, timeoutCancel := context.WithTimeout(signalCtx, *timeout)
	defer timeoutCancel()

	me, err := me()
	if err != nil {
		log.Fatal("Error getting user:", err)
	}

	eg, egctx := errgroup.WithContext(timeoutCtx)

	for repo := range ghreposync.Filter(cfg.WorkDir, cfg.Repositories) {
		eg.Go(func() error {
			splits := strings.Split(repo, "/")
			repoName := splits[len(splits)-1]
			slog.Info("Start to sync the repository", slog.String("name", repo))
			if err := run(egctx, me, repoName); err != nil {
				slog.Error("Failed to sync the repository", slog.String("name", repo), slog.String("error", err.Error()))
				return err
			}
			slog.Info("Successfully synced the repository", slog.String("name", repo))
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		log.Fatal("Error syncing repos:", err)
	}
}

func run(ctx context.Context, username, repo string) error {
	entry := "gh repo sync %s/%s"
	cmd := exec.CommandContext(ctx, "sh", "-c", fmt.Sprintf(entry, username, repo))
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	cmd.Env = append(os.Environ(), "GH_NO_PROMPT=1")
	cmd.Cancel = func() error {
		return cmd.Process.Signal(os.Interrupt)
	}
	cmd.WaitDelay = 5 * time.Second
	return cmd.Run()
}

func me() (string, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return "", err
	}
	response := struct{ Login string }{}
	err = client.Get("user", &response)
	if err != nil {
		return "", err
	}
	return response.Login, nil
}
