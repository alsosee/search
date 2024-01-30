package main

import (
	"bufio"
	"context"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/meilisearch/meilisearch-go"
	gitignore "github.com/sabhiram/go-gitignore"
	"gopkg.in/yaml.v3"

	"github.com/alsosee/finder/structs"
)

// Indexer reads files and writes them to a MeiliSearch index.
type Indexer struct {
	client *meilisearch.Client
	ctx    context.Context

	state   map[string]string
	muState sync.Mutex

	infoDir  string
	mediaDir string
}

// NewIndexer creates a new Indexer.
func NewIndexer(
	client *meilisearch.Client,
	infoDir string,
	mediaDir string,
) *Indexer {
	return &Indexer{
		client:   client,
		ctx:      context.Background(),
		state:    make(map[string]string),
		infoDir:  infoDir,
		mediaDir: mediaDir,
	}
}

// Index reads files from the info directory and writes them to the MeiliSearch.
func (i *Indexer) Index(stateFile, ignoreFile, index string) error {
	ignore, err := processIgnoreFile(filepath.Join(i.infoDir, ignoreFile))
	if err != nil {
		return fmt.Errorf("processing ignore file: %w", err)
	}

	absDir, err := filepath.Abs(i.infoDir)
	if err != nil {
		return fmt.Errorf("getting absolute path: %w", err)
	}
	err = filepath.Walk(absDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walking directory: %w", err)
		}

		relPath := strings.TrimPrefix(path, absDir+string(filepath.Separator))

		if ignore.MatchesPath(relPath) {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if ext != ".yml" && ext != ".yaml" { // && ext != ".md" {
			return nil
		}

		if err := i.addFile(path, relPath); err != nil {
			return fmt.Errorf("adding file: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("walking directory: %w", err)
	}

	state, err := readStateFromFile(stateFile)
	if err != nil {
		return fmt.Errorf("reading state file: %w", err)
	}

	if err := i.updateIndex(state, index); err != nil {
		return fmt.Errorf("updating index: %w", err)
	}

	if err := writeStateToFile(stateFile, i.state); err != nil {
		return fmt.Errorf("writing state file: %w", err)
	}

	return nil
}

func (i *Indexer) addFile(path, relPath string) error {
	i.muState.Lock()
	defer i.muState.Unlock()

	hash, err := hash(path)
	if err != nil {
		return fmt.Errorf("hashing file: %w", err)
	}

	i.state[relPath] = hash
	return nil
}

func (i *Indexer) updateIndex(oldState map[string]string, index string) error {
	// find deleted files
	var toDelete []string
	for path := range oldState {
		if _, ok := i.state[path]; !ok {
			toDelete = append(toDelete, path)
		}
	}

	// find new and changed files
	var toUpdate []string
	for path, hash := range i.state {
		if oldHash, ok := oldState[path]; !ok || oldHash != hash {
			toUpdate = append(toUpdate, path)
		}
	}

	if err := i.deleteFromIndex(toDelete, index); err != nil {
		return fmt.Errorf("deleting documents: %w", err)
	}

	if err := i.addToIndex(toUpdate, index); err != nil {
		return fmt.Errorf("adding documents: %w", err)
	}

	return nil
}

func (i *Indexer) deleteFromIndex(paths []string, index string) error {
	if len(paths) == 0 {
		return nil
	}

	// todo fix IDs
	task, err := i.client.Index(index).DeleteDocuments(paths)
	if err != nil {
		return err
	}

	err = i.waitForTask(task.TaskUID, time.Minute*2)
	if err != nil {
		return fmt.Errorf("waiting for task %q: %w", task.TaskUID, err)
	}

	return nil
}

func (i *Indexer) addToIndex(paths []string, index string) error {
	documents := []*structs.Content{}

	for _, path := range paths {
		document, err := i.processFile(path)
		if err != nil {
			return fmt.Errorf("processing file %q: %w", path, err)
		}
		documents = append(documents, document)
	}

	if len(documents) == 0 {
		log.Printf("No documents to add to index %q", index)
		return nil
	}

	log.Printf("Adding %d documents to index %q", len(documents), index)
	for _, document := range documents {
		log.Printf("  %s", document.Source)
	}

	tasks, err := i.client.Index(index).AddDocumentsInBatches(documents, 100, "ID")
	if err != nil {
		return err
	}

	for _, task := range tasks {
		err := i.waitForTask(task.TaskUID, time.Minute*2)
		if err != nil {
			return fmt.Errorf("waiting for task %q: %w", task.TaskUID, err)
		}
	}
	return err
}

func (i *Indexer) processFile(path string) (*structs.Content, error) {
	switch filepath.Ext(path) {
	case ".yml", ".yaml":
		return i.processYAMLFile(path)
	// case ".md":
	// 	return i.processMarkdownFile(path)
	default:
		return nil, fmt.Errorf("unknown file type: %s", path)
	}
}

func (i *Indexer) processYAMLFile(path string) (*structs.Content, error) {
	file, err := os.Open(filepath.Join(i.infoDir, path))
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	var content structs.Content
	if err := yaml.NewDecoder(file).Decode(&content); err != nil {
		return nil, fmt.Errorf("decoding file: %w", err)
	}

	content.Source = path

	id := removeFileExtention(path)

	content.Source = path
	content.ID = formatID(id)
	content.Image = i.getImageForPath(id)

	// add image to Characters
	for _, character := range content.Characters {
		character.Image = i.getImageForPath(filepath.Join(id, "Characters", character.Name))
		if character.Actor != "" {
			character.ActorImage = i.getImageForPath("People/" + character.Actor)
		}
	}

	return &content, nil
}

func (i *Indexer) getImageForPath(path string) *structs.Media {
	dir := filepath.Dir(path)
	if dir == "." {
		dir = ""
	}

	mediaAbsPath, err := filepath.Abs(i.mediaDir)
	if err != nil {
		log.Printf("Error getting absolute path for %q: %v", i.mediaDir, err)
		return nil
	}

	// read .thumb.yml file in media directory
	thumbFile := filepath.Join(mediaAbsPath, dir, ".thumbs.yml")
	if _, err := os.Stat(thumbFile); os.IsNotExist(err) {
		return nil
	}

	media, err := structs.ParseMediaFile(thumbFile)
	if err != nil {
		log.Printf("Error parsing media file %q: %v", thumbFile, err)
		return nil
	}

	if len(media) == 0 {
		return nil
	}

	for _, m := range media {
		mediaPath := filepath.Join(dir, removeFileExtention(m.Path))
		if mediaPath == path {
			return &m
		}
	}

	return nil
}

func (i *Indexer) waitForTask(taskID int64, timeout time.Duration) error {
	// increase default context timeout from 5s to 2m to wait for slow
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	_, err := i.client.WaitForTask(
		taskID,
		meilisearch.WaitParams{
			Context:  ctx,
			Interval: time.Millisecond * 500,
		},
	)
	if err != nil {
		return err
	}

	// check task status
	task, err := i.client.GetTask(taskID)
	if err != nil {
		return err
	}

	if task.Status == "failed" {
		return fmt.Errorf("task failed: %s", task.Error)
	}
	return nil
}

func hash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	hash := crc32.NewIEEE()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("calculating CRC32 checksum: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum32()), nil
}

func readStateFromFile(stateFile string) (map[string]string, error) {
	state := make(map[string]string)

	f, err := os.Open(stateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return state, nil
		}
		return nil, fmt.Errorf("opening state file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid state file format")
		}
		relPath := parts[0]
		hash := parts[1]
		state[relPath] = hash
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading state file: %w", err)
	}

	return state, nil
}

func writeStateToFile(stateFile string, state map[string]string) error {
	f, err := os.Create(stateFile)
	if err != nil {
		return fmt.Errorf("creating state file: %w", err)
	}
	defer f.Close()

	stateSlice := make([]string, 0, len(state))

	for relPath, hash := range state {
		stateSlice = append(stateSlice, fmt.Sprintf("%s\t%s", relPath, hash))
	}

	sort.Strings(stateSlice)

	for _, line := range stateSlice {
		if _, err := f.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("writing state file: %w", err)
		}
	}

	return nil
}

func processIgnoreFile(path string) (*gitignore.GitIgnore, error) {
	ignore := &gitignore.GitIgnore{}
	if _, err := os.Stat(path); err == nil {
		ignore, err = gitignore.CompileIgnoreFile(path)
		if err != nil {
			return nil, fmt.Errorf("compiling ignore file: %w", err)
		}
	} else {
		log.Printf("Ignore file %q not found, ignoring", path)
	}

	return ignore, nil
}

func removeFileExtention(path string) string {
	withoutExt := path[:len(path)-len(filepath.Ext(path))]
	if withoutExt != "" {
		return withoutExt
	}
	return path
}

var reNonID = regexp.MustCompile("[^a-zA-Z0-9-_]")

// formatID formats an ID for MeiliSearch.
// A document identifier can be of type integer or string,
// only composed of alphanumeric characters (a-z A-Z 0-9), hyphens (-) and underscores (_).
func formatID(id string) string {
	return reNonID.ReplaceAllString(id, "_")
}
