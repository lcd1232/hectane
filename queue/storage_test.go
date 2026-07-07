package queue

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	var (
		data        = []byte("test")
		numMessages = 5
	)
	d, err := ioutil.TempDir(os.TempDir(), "")
	require.NoError(t, err)
	defer os.RemoveAll(d)
	s := NewStorage(d, 0)
	w, body, err := s.NewBody()
	require.NoError(t, err)
	_, err = w.Write(data)
	require.NoError(t, err)
	require.NoError(t, w.Close())
	for i := 0; i < numMessages; i++ {
		require.NoError(t, s.SaveMessage(&Message{}, body))
	}
	messages, err := s.LoadMessages()
	require.NoError(t, err)
	for _, m := range messages {
		r, err := s.GetMessageBody(m)
		require.NoError(t, err)
		b, err := ioutil.ReadAll(r)
		require.NoError(t, err)
		assert.Equal(t, data, b)
		require.NoError(t, r.Close())
		require.NoError(t, s.DeleteMessage(m))
	}
	e, err := ioutil.ReadDir(s.directory)
	require.NoError(t, err)
	if len(e) != 0 {
		t.Fatalf("%d != 0", len(e))
	}
}

func TestStorageTTL(t *testing.T) {
	var (
		data        = []byte("test")
		numMessages = 5
	)
	d, err := ioutil.TempDir(os.TempDir(), "")
	require.NoError(t, err)
	defer os.RemoveAll(d)
	s := NewStorage(d, 10*time.Second)
	w, body, err := s.NewBody()
	require.NoError(t, err)
	_, err = w.Write(data)
	require.NoError(t, err)
	require.NoError(t, w.Close())
	for i := 0; i < numMessages; i++ {
		require.NoError(t, s.SaveMessage(&Message{}, body))
	}
	messages, err := s.LoadMessages()
	require.NoError(t, err)
	assert.Len(t, messages, numMessages)
	time.Sleep(10 * time.Second)
	messages, err = s.LoadMessages()
	require.NoError(t, err)
	assert.Len(t, messages, 0)
	// Once every message has expired the body directory must be removed too,
	// otherwise empty directories accumulate in the spool.
	entries, err := ioutil.ReadDir(s.directory)
	require.NoError(t, err)
	assert.Empty(t, entries, "expired body directory should be removed")
}

// TestStorageDeleteMessageEmptyDirectory covers the case where the body file is
// already gone (e.g. removed by TTL expiry) when the last message is deleted.
// DeleteMessage must remove the now-empty directory and must not return an
// error for the empty-directory case.
func TestStorageDeleteMessageEmptyDirectory(t *testing.T) {
	d, err := ioutil.TempDir(os.TempDir(), "")
	require.NoError(t, err)
	defer os.RemoveAll(d)
	s := NewStorage(d, 0)
	w, body, err := s.NewBody()
	require.NoError(t, err)
	require.NoError(t, w.Close())
	m := &Message{}
	require.NoError(t, s.SaveMessage(m, body))
	// Simulate the body file having been removed independently, leaving only
	// the message file behind.
	require.NoError(t, os.Remove(s.bodyFilename(body)))

	require.NoError(t, s.DeleteMessage(m))

	entries, err := ioutil.ReadDir(s.directory)
	require.NoError(t, err)
	assert.Empty(t, entries, "body directory should be removed once empty")
}

// TestStorageDeleteMessageKeepsPendingSibling verifies that a body directory
// shared by several messages is only removed once the last message is deleted.
func TestStorageDeleteMessageKeepsPendingSibling(t *testing.T) {
	d, err := ioutil.TempDir(os.TempDir(), "")
	require.NoError(t, err)
	defer os.RemoveAll(d)
	s := NewStorage(d, 0)
	w, body, err := s.NewBody()
	require.NoError(t, err)
	require.NoError(t, w.Close())
	m1, m2 := &Message{}, &Message{}
	require.NoError(t, s.SaveMessage(m1, body))
	require.NoError(t, s.SaveMessage(m2, body))

	// Deleting the first message must keep the directory for the second.
	require.NoError(t, s.DeleteMessage(m1))
	_, err = os.Stat(s.bodyDirectory(body))
	require.NoError(t, err, "directory must survive while a message is pending")

	// Deleting the last message removes the directory.
	require.NoError(t, s.DeleteMessage(m2))
	entries, err := ioutil.ReadDir(s.directory)
	require.NoError(t, err)
	assert.Empty(t, entries, "body directory should be removed with no messages left")
}

// TestStorageOrphanedDirectoryCleanup verifies that a directory without a body
// file (e.g. left by a crashed NewBody) is removed by LoadMessages.
func TestStorageOrphanedDirectoryCleanup(t *testing.T) {
	d, err := ioutil.TempDir(os.TempDir(), "")
	require.NoError(t, err)
	defer os.RemoveAll(d)
	s := NewStorage(d, 0)
	require.NoError(t, os.MkdirAll(s.bodyDirectory("orphaned-body-id"), 0700))

	_, err = s.LoadMessages()
	require.NoError(t, err)

	entries, err := ioutil.ReadDir(s.directory)
	require.NoError(t, err)
	assert.Empty(t, entries, "orphaned directory should be removed")
}
