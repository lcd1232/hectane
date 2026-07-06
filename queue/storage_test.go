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
	// Verify that the body directory was removed after all messages expired.
	e, err := ioutil.ReadDir(s.directory)
	require.NoError(t, err)
	assert.Len(t, e, 0)
}

// TestStorageOrphanedDirectoryCleanup verifies that directories without a
// body file (e.g. left by a crashed NewBody call) are removed by LoadMessages.
func TestStorageOrphanedDirectoryCleanup(t *testing.T) {
	d, err := ioutil.TempDir(os.TempDir(), "")
	require.NoError(t, err)
	defer os.RemoveAll(d)
	s := NewStorage(d, 0)
	// Simulate a partial NewBody failure by manually creating an empty directory.
	orphan := s.bodyDirectory("orphaned-body-id")
	require.NoError(t, os.MkdirAll(orphan, 0700))
	_, err = s.LoadMessages()
	require.NoError(t, err)
	// The orphaned directory should have been removed.
	e, err := ioutil.ReadDir(s.directory)
	require.NoError(t, err)
	assert.Len(t, e, 0)
}
