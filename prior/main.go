package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {

	pretty := flag.Bool("pretty", false, "pretty output")
	flag.Parse()

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	var w io.Writer = os.Stdout

	if *pretty {
		w = tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
	}

	var all []*github.Repository

	for page := 1; ; page++ {
		repos, _, err := client.Repositories.List(ctx, "", &github.RepositoryListOptions{ListOptions: github.ListOptions{Page: page}})
		if err != nil {
			log.Printf("error fetching repositories page %d = %+v\n", page, err)
			break
		}

		log.Printf("fetched %d results", len(repos))

		if len(repos) == 0 {
			log.Printf("done")
			break
		}

		all = append(all, repos...)
	}

	log.Printf("total %d repos", len(all))

	sort.Sort(ByCreatedAt(all))

	fmt.Fprintln(w, strings.Join([]string{
		"Name",
		"Description",
		"Created",
		"Last Modification"}, "\t"))

	for _, r := range all {
		// ignore forks
		if r.GetFork() {
			log.Printf("ignoring fork %q", r.GetFullName())
			continue
		}

		if r.GetDescription() == "" {
			log.Printf("missing description for %q", r.GetFullName())
		}

		fmt.Fprintln(w, strings.Join([]string{
			r.GetFullName(),
			r.GetDescription(),
			r.GetCreatedAt().Format("2006-01-02"),
			r.GetPushedAt().Format("2006-01-02")}, "\t"))
	}

	if f, ok := w.(interface {
		Flush() error
	}); ok {
		f.Flush()
	}
}

type ByCreatedAt []*github.Repository

func (b ByCreatedAt) Len() int           { return len(b) }
func (b ByCreatedAt) Less(i, j int) bool { return b[i].GetCreatedAt().Before(b[j].GetCreatedAt().Time) }
func (b ByCreatedAt) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
