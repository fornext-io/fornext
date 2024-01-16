package executor

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/cockroachdb/pebble/vfs"
)

type Storage struct {
	db *pebble.DB
}

func NewStorage(dir string) (*Storage, error) {
	pbOptions := &pebble.Options{
		BytesPerSync:          512 * 1024,
		Cache:                 pebble.NewCache(8 * 1024 * 1024),
		Cleaner:               pebble.DeleteCleaner{},
		Comparer:              pebble.DefaultComparer,
		DebugCheck:            nil,
		DisableWAL:            false,
		ErrorIfExists:         false,
		ErrorIfNotPristine:    false,
		EventListener:         &pebble.EventListener{},
		Filters:               map[string]pebble.FilterPolicy{},
		FlushDelayDeleteRange: time.Duration(0),
		FlushDelayRangeKey:    time.Duration(0),
		FlushSplitBytes:       0,
		FormatMajorVersion:    pebble.FormatNewest,
		FS:                    vfs.Default,
		Lock:                  nil,
		Levels: []pebble.LevelOptions{
			{
				BlockSize:      64 * 1024,
				Compression:    pebble.NoCompression,
				TargetFileSize: 32 * 1024 * 1024,
			}, {
				BlockSize:      64 * 1024,
				Compression:    pebble.ZstdCompression,
				TargetFileSize: 64 * 1024 * 1024,
			},
		},
		LoggerAndTracer: nil,
		MaxOpenFiles:    1000,
	}
	db, err := pebble.Open(dir, pbOptions)
	if err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Set(ctx context.Context, key []byte, value []byte) error {
	batch := s.db.NewIndexedBatch()
	defer func() {
		err := batch.Close()
		if err != nil {
			panic(err)
		}
	}()

	err := batch.Set(key, value, pebble.Sync)
	if err != nil {
		return err
	}

	err = batch.Commit(pebble.Sync)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Get(ctx context.Context, key []byte) (value []byte, err error) {
	batch := s.db.NewIndexedBatch()
	defer func() {
		err := batch.Close()
		if err != nil {
			panic(err)
		}
	}()

	v, closer, err := batch.Get(key)
	defer func() {
		err := closer.Close()
		if err != nil {
			panic(err)
		}
	}()

	value = make([]byte, len(v))
	copy(value, v)
	return value, nil
}

func Set[T any](ctx context.Context, s *Storage, key string, value T) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.Set(ctx, []byte(key), data)
}

func Get[T any](ctx context.Context, s *Storage, key string) (value *T, err error) {
	data, err := s.Get(ctx, []byte(key))
	if err != nil {
		return value, err
	}

	err = json.Unmarshal(data, &value)
	if err != nil {
		return value, err
	}

	return value, nil
}
