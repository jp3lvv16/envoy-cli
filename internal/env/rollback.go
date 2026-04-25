package env

import "fmt"

// RollbackResult holds the outcome of a rollback operation.
type RollbackResult struct {
	Set     *Set
	Version int
}

// Rollback reverts a Set to a previous snapshot captured in a History.
// version is 0-based; 0 means the most recent snapshot, 1 means one before that, etc.
func Rollback(h *History, version int) (*RollbackResult, error) {
	if h == nil {
		return nil, fmt.Errorf("rollback: history must not be nil")
	}
	if version < 0 {
		return nil, fmt.Errorf("rollback: version must be >= 0, got %d", version)
	}

	snaps := h.Snapshots()
	if len(snaps) == 0 {
		return nil, fmt.Errorf("rollback: no snapshots available")
	}

	// Most recent snapshot is last in the slice.
	idx := len(snaps) - 1 - version
	if idx < 0 {
		return nil, fmt.Errorf("rollback: version %d out of range (have %d snapshot(s))", version, len(snaps))
	}

	restored, err := RestoreSnapshot(snaps[idx])
	if err != nil {
		return nil, fmt.Errorf("rollback: %w", err)
	}

	return &RollbackResult{
		Set:     restored,
		Version: version,
	}, nil
}

// RollbackToLatest is a convenience wrapper that rolls back to the most
// recent snapshot (version 0).
func RollbackToLatest(h *History) (*RollbackResult, error) {
	return Rollback(h, 0)
}
