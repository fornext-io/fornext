package executor

import (
	"context"
)

type LogStream struct {
	ev chan interface{}
}

func (s *LogStream) Append(ctx context.Context, v interface{}) error {
	s.ev <- v
	return nil
}

func (s *LogStream) Read(ctx context.Context) <-chan interface{} {
	return s.ev
}
