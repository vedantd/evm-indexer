package planner

import (
	"context"
	"fmt"
)

// HeadProvider abstracts how we learn the current head number.
// We'll fake this in tests and later implement an RPC version.
type HeadProvider interface {
	HeadNumber(ctx context.Context) (uint64, error)
}

// Planner emits historical block numbers in batches until it reaches (head - safetyWindow).
// It pushes numbers into the provided Out channel; backpressure is handled by the channel's buffering/consumer.
type Planner struct {
	Heads        HeadProvider
	BatchSize    uint64 // e.g., 100
	SafetyWindow uint64 // e.g., 6
}

// Plan emits block numbers starting from `from` up to target = max(0, head - safetyWindow).
// It blocks on `Out` writes if the channel is full; cancel via ctx to stop early.
func (p *Planner) Plan(ctx context.Context, from uint64, Out chan<- uint64) error {
	if p.Heads == nil {
		return fmt.Errorf("planner: Heads (HeadProvider) is nil")
	}
	if p.BatchSize == 0 {
		p.BatchSize = 100
	}

	// Get the head at the time of this planning pass.
	head, err := p.Heads.HeadNumber(ctx)
	if err != nil {
		return fmt.Errorf("planner: head: %w", err)
	}

	target := uint64(0)
	if head > p.SafetyWindow {
		target = head - p.SafetyWindow
	}
	if from > target {
		// Nothing to do.
		return nil
	}

	next := from
	for {
		// Compute range [start..end] for this batch, capped at target
		end := next + p.BatchSize - 1
		if end > target {
			end = target
		}

		// Emit the batch range
		for n := next; n <= end; n++ {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case Out <- n:
			}
		}

		// Advance or finish
		if end == target {
			return nil
		}
		next = end + 1
	}
}
