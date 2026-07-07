package queue

import (
	"time"

	"github.com/pborman/uuid"

	"encoding/json"
	"io"
	"os"
	"path"
	"strings"
	"sync"
)

const (
	bodyFilename     = "body"
	messageExtension = ".message"
)

// Message metadata.
type Message struct {
	id      string
	body    string
	Host    string   `json:"host"`
	From    string   `json:"from"`
	To      []string `json:"to"`
	Attempt int      `json:"attempt"`
	// Metadata can store any data useful for custom processors.
	Metadata map[string]string `json:"metadata"`
}

// Manager for message metadata and body on disk. All methods are safe to call
// from multiple goroutines.
type Storage struct {
	m         sync.Mutex
	directory string
	ttl       time.Duration
}

// Determine the path to the directory containing the specified body.
func (s *Storage) bodyDirectory(body string) string {
	return path.Join(s.directory, body)
}

// Determine the filename of the specified body.
func (s *Storage) bodyFilename(body string) string {
	return path.Join(s.bodyDirectory(body), bodyFilename)
}

// Determine the filename of the specified message.
func (s *Storage) messageFilename(m *Message) string {
	return path.Join(s.directory, m.body, m.id) + messageExtension
}

// Load all messages with the specified body.
func (s *Storage) loadMessages(body string) []*Message {
	messages := make([]*Message, 0, 1)
	if files, err := os.ReadDir(s.bodyDirectory(body)); err == nil {
		for _, f := range files {
			info, err := f.Info()
			if err != nil {
				continue // File vanished between listing and stat; skip it.
			}
			if s.ttl > 0 && time.Since(info.ModTime()) > s.ttl {
				// If the file is older than the TTL, delete it.
				if err := os.Remove(path.Join(s.bodyDirectory(body), f.Name())); err != nil {
					continue // Ignore errors, just skip this file.
				}
				continue // Skip this file, it is too old.
			}
			if strings.HasSuffix(f.Name(), messageExtension) {
				m := &Message{
					id:   strings.TrimSuffix(f.Name(), messageExtension),
					body: body,
				}
				if r, err := os.Open(s.messageFilename(m)); err == nil {
					if err := json.NewDecoder(r).Decode(m); err == nil {
						messages = append(messages, m)
					}
					r.Close()
				}
			}
		}
	}
	// Once no message files remain (e.g. they all expired) drop the whole body
	// directory so empty directories don't pile up in the spool.
	_ = s.removeBodyDirectoryIfEmpty(body)
	return messages
}

// removeBodyDirectoryIfEmpty removes the body directory when it no longer holds
// any message files. The body file itself does not count, so the directory is
// cleaned up once every message sharing it has been delivered or expired. A
// message file that is still on disk (even one that failed to load) keeps the
// directory, so nothing pending is silently discarded here.
func (s *Storage) removeBodyDirectoryIfEmpty(body string) error {
	dir := s.bodyDirectory(body)
	files, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), messageExtension) {
			return nil
		}
	}
	return os.RemoveAll(dir)
}

// NewStorage creates a Storage instance for the specified directory.
// TTL is the time to live for messages, after which they will be deleted.
func NewStorage(directory string, ttl time.Duration) *Storage {
	return &Storage{
		directory: directory,
		ttl:       ttl,
	}
}

// Create a new message body. The writer must be closed after writing the
// message body.
func (s *Storage) NewBody() (io.WriteCloser, string, error) {
	body := uuid.New()
	if err := os.MkdirAll(s.bodyDirectory(body), 0700); err != nil {
		return nil, "", err
	}
	w, err := os.OpenFile(s.bodyFilename(body), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		// MkdirAll already created the directory; don't leak it on failure.
		_ = os.RemoveAll(s.bodyDirectory(body))
		return nil, "", err
	}
	return w, body, nil
}

// Load messages from the storage directory. Any messages that could not be
// loaded are ignored.
func (s *Storage) LoadMessages() ([]*Message, error) {
	directories, err := os.ReadDir(s.directory)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		return []*Message{}, nil
	}
	var messages []*Message
	for _, d := range directories {
		if !d.IsDir() {
			continue
		}
		_, err := os.Stat(s.bodyFilename(d.Name()))
		switch {
		case err == nil:
			messages = append(messages, s.loadMessages(d.Name())...)
		case os.IsNotExist(err):
			// The directory has no body file (e.g. a crashed NewBody). It can
			// never be delivered, so remove it instead of scanning it forever.
			_ = os.RemoveAll(s.bodyDirectory(d.Name()))
		default:
			// Any other error (I/O, permissions) is transient; leave the
			// directory untouched rather than risk deleting valid messages.
		}
	}
	return messages, nil
}

// Save the specified message to disk.
func (s *Storage) SaveMessage(m *Message, body string) error {
	s.m.Lock()
	defer s.m.Unlock()
	if m.id == "" {
		m.id = uuid.New()
	}
	m.body = body
	w, err := os.OpenFile(s.messageFilename(m), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer w.Close()
	if err := json.NewEncoder(w).Encode(m); err != nil {
		return err
	}
	return nil
}

// Retrieve a reader for the message body.
func (s *Storage) GetMessageBody(m *Message) (io.ReadCloser, error) {
	s.m.Lock()
	defer s.m.Unlock()
	return os.Open(s.bodyFilename(m.body))
}

// Delete the specified message. The message body is also deleted if no more
// messages exist.
func (s *Storage) DeleteMessage(m *Message) error {
	s.m.Lock()
	defer s.m.Unlock()
	if err := os.Remove(s.messageFilename(m)); err != nil {
		return err
	}
	return s.removeBodyDirectoryIfEmpty(m.body)
}
