// Search app works as a GitHub action that checks state file for any changes
// and updates Meilisearch index accordingly.
package main

import (
	"log"
	"strings"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/meilisearch/meilisearch-go"
)

type config struct {
	InfoDirectory  string        `env:"INPUT_INFO" short:"i" long:"info" description:"Directory that contains info files" default:"info"`
	MediaDirectory string        `env:"INPUT_MEDIA" short:"m" long:"media" description:"Directory that contains media files" default:"media"`
	StateFile      string        `env:"INPUT_STATE_FILE" long:"state-file" description:"path to state file" default:".state"`
	IgnoreFile     string        `env:"INPUT_IGNORE_FILE" short:"g" long:"ignore" description:"Path to .gitignore file" default:".ignore"`
	Host           string        `env:"INPUT_HOST" long:"host" description:"search host" default:"http://127.0.0.1:7700/"`
	IndexName      string        `env:"INPUT_INDEX" long:"index" description:"search index name" default:"info"`
	MasterKey      string        `env:"INPUT_MASTER_KEY" long:"master-key" description:"search master key"`
	Timeout        time.Duration `env:"INPUT_TIMEOUT" long:"timeout" description:"search timeout" default:"5s"`
}

var cfg config

func main() {
	if _, err := flags.Parse(&cfg); err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:    cfg.Host,
		APIKey:  strings.TrimSpace(cfg.MasterKey),
		Timeout: cfg.Timeout,
	})

	indexer := NewIndexer(
		client,
		cfg.InfoDirectory,
		cfg.MediaDirectory,
	)
	if err := indexer.Index(
		cfg.StateFile,
		cfg.IgnoreFile,
		cfg.IndexName,
	); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
