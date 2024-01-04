package executor

import (
	"strconv"
	"sync"
	"time"
)

// Timestamp is the represent of `hybrid logical clock`
type Timestamp struct {
	// WallTime is the physical unix epoch time expressed in
	// seconds.
	WallTime uint32

	// Logical is an sequential clock to captures causality for events
	// whose wall times are equal.
	Logical uint32
}

func (t Timestamp) AsUint64() uint64 {
	return uint64(t.Logical)<<32 + uint64(t.WallTime)
}

var (
	chars = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func (t Timestamp) AsString() string {
	// buf := make([]byte, 8)
	// binary.BigEndian.PutUint64(buf, t.AsUint64())
	// return string(buf)
	return strconv.FormatUint(t.AsUint64(), 16)
	// v := t.AsUint64()
	// buf := bytes.Buffer{}
	// buf.Grow(12)

	// for v > 0 {
	// 	vv := v % 62
	// 	v /= 62
	// 	buf.WriteByte(chars[vv])
	// }
	// return buf.String()
}

func ParseTimestamp(t uint64) Timestamp {
	return Timestamp{
		Logical:  uint32(t >> 32),
		WallTime: uint32(t),
	}
}

type hybridLogicalClock struct {
	mu sync.Mutex

	// Timestamp is the current `hybrid logical clock` value of current node.
	// It always been the max value of this node had seen.
	timestamp Timestamp
}

func newHybridLogicalClock() *hybridLogicalClock {
	return &hybridLogicalClock{}
}

func (c *hybridLogicalClock) Update(t Timestamp) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// We have an bigger Timestamp than remote nodes, so we just keep our's
	// value.
	if t.WallTime < c.timestamp.WallTime {
		return
	}

	if t.WallTime > c.timestamp.WallTime {
		// The remote clock is ahead of ours, so we update
		// our Timestamp using remote's value.
		c.timestamp.WallTime = t.WallTime
		c.timestamp.Logical = t.Logical
	} else if t.WallTime == c.timestamp.WallTime {
		// If remote clock's WallTime is equal with our's, but logical clock
		// is ahead of ours, use it.
		if t.Logical > c.timestamp.Logical {
			c.timestamp.Logical = t.Logical
		}
	}
}

func (c *hybridLogicalClock) Next() Timestamp {
	now := uint32(time.Now().Unix())

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.timestamp.WallTime >= now {
		// The current WallTime is ahead of current node's physical time,
		// so we only tick the logical clock.
		// This could happened in below scenario:
		// 1. current node's physical time non-monotonic updates
		// 2. receive an Timestamp with bigger WallTime from other nodes
		c.timestamp.Logical++
	} else {
		// Use physical time as WallTime, and reset the logical clock to zero.
		c.timestamp.WallTime = now
		c.timestamp.Logical = 0
	}

	return c.timestamp
}
