package mem

import (
	"context"
	"fmt"
	"sync"
	"time"

	apierrors "github.com/gomods/athens/pkg/errors"
	"github.com/gomods/athens/pkg/index"
)

// New returns a new in-memory indexer.
func New() index.Indexer {
	return &indexer{}
}

type indexer struct {
	mu    sync.RWMutex
	lines []*index.Line
}

func (i *indexer) Index(_ context.Context, mod, ver string) error {
	const op apierrors.Op = "mem.Index"
	i.mu.Lock()
	defer i.mu.Unlock()
	for _, l := range i.lines {
		if l.Path == mod && l.Version == ver {
			return apierrors.E(op, fmt.Sprintf("%s@%s already indexed", mod, ver), apierrors.KindAlreadyExists)
		}
	}
	i.lines = append(i.lines, &index.Line{
		Path:      mod,
		Version:   ver,
		Timestamp: time.Now(),
	})
	return nil
}

func (i *indexer) Lines(_ context.Context, since time.Time, limit int) ([]*index.Line, error) {
	lines := []*index.Line{}
	var count int
	i.mu.RLock()
	defer i.mu.RUnlock()
	for _, line := range i.lines {
		if count >= limit {
			break
		}
		if since.After(line.Timestamp) {
			continue
		}
		lines = append(lines, line)
		count++
	}
	return lines, nil
}
