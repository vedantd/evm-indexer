package planner

import (
	"context"
	"testing"
	"time"
)

// fake head provider for tests
type fakeHeads uint64

func (f fakeHeads) HeadNumber(ctx context.Context) (uint64, error) {
	return uint64(f), nil
}

func collect(out <-chan uint64) []uint64 {
	var xs []uint64
	for {
		select {
		case n, ok := <-out:
			if !ok {
				return xs
			}
			xs = append(xs, n)
		default:
			return xs
		}
	}
}

func TestPlan_BasicRange(t *testing.T) {
	ctx := context.Background()
	p := &Planner{
		Heads:        fakeHeads(1000),
		BatchSize:    100,
		SafetyWindow: 6,
	}

	// target = 1000 - 6 = 994
	out := make(chan uint64, 3000) // large buffer to avoid blocking in test

	if err := p.Plan(ctx, 800, out); err != nil {
		t.Fatalf("plan error: %v", err)
	}
	got := collect(out)

	// Expect 800..994 inclusive
	if len(got) != (994 - 800 + 1) {
		t.Fatalf("unexpected length: got=%d want=%d", len(got), (994 - 800 + 1))
	}
	if got[0] != 800 || got[len(got)-1] != 994 {
		t.Fatalf("unexpected bounds: first=%d last=%d", got[0], got[len(got)-1])
	}
}

func TestPlan_PartialLastBatch(t *testing.T) {
	ctx := context.Background()
	p := &Planner{
		Heads:        fakeHeads(250),
		BatchSize:    80,
		SafetyWindow: 10,
	}
	// target = 240
	// from=100 → batches: [100..179], [180..239], [240..240]
	out := make(chan uint64, 1000)
	if err := p.Plan(ctx, 100, out); err != nil {
		t.Fatalf("plan error: %v", err)
	}
	got := collect(out)
	if got[0] != 100 || got[len(got)-1] != 240 {
		t.Fatalf("unexpected bounds: first=%d last=%d", got[0], got[len(got)-1])
	}
	// quick spot checks at batch boundaries
	if got[79] != 179 || got[80] != 180 {
		t.Fatalf("batch boundary wrong: got[79]=%d got[80]=%d", got[79], got[80])
	}
}

func TestPlan_NoWorkWhenFromBeyondTarget(t *testing.T) {
	ctx := context.Background()
	p := &Planner{
		Heads:        fakeHeads(100),
		BatchSize:    50,
		SafetyWindow: 10,
	}
	// target = 90, from=200 => nothing
	out := make(chan uint64, 1)
	if err := p.Plan(ctx, 200, out); err != nil {
		t.Fatalf("plan error: %v", err)
	}
	got := collect(out)
	if len(got) != 0 {
		t.Fatalf("expected no numbers, got %v", got)
	}
}

func TestPlan_ContextCancelStopsEmission(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	p := &Planner{
		Heads:        fakeHeads(1_000_000),
		BatchSize:    100_000,
		SafetyWindow: 0,
	}
	out := make(chan uint64, 10) // tiny buffer so it blocks quickly

	// Run Plan in a goroutine; cancel soon after to ensure it returns
	errCh := make(chan error, 1)
	go func() {
		errCh <- p.Plan(ctx, 0, out)
	}()

	time.Sleep(20 * time.Millisecond)
	cancel()

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatalf("expected context error, got nil")
		}
	case <-time.After(1 * time.Second):
		t.Fatalf("plan did not return after cancel")
	}
}
