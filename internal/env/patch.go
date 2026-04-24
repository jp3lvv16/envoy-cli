package env

import "fmt"

// PatchOp represents a single patch operation on an env set.
type PatchOp struct {
	Op    string // "set", "delete", "rename"
	Key   string
	Value string // used by "set"
	NewKey string // used by "rename"
}

// Patch applies a slice of PatchOps to the given Set in order.
// Supported operations:
//   - "set":    sets Key to Value
//   - "delete": removes Key
//   - "rename": renames Key to NewKey, preserving the value
//
// Returns an error if the set is nil, an op is unknown, or an op fails.
func Patch(s *Set, ops []PatchOp) error {
	if s == nil {
		return fmt.Errorf("patch: set must not be nil")
	}
	for i, op := range ops {
		switch op.Op {
		case "set":
			if err := s.Put(op.Key, op.Value); err != nil {
				return fmt.Errorf("patch[%d] set %q: %w", i, op.Key, err)
			}
		case "delete":
			if err := s.Delete(op.Key); err != nil {
				return fmt.Errorf("patch[%d] delete %q: %w", i, op.Key, err)
			}
		case "rename":
			val, err := s.Get(op.Key)
			if err != nil {
				return fmt.Errorf("patch[%d] rename %q: %w", i, op.Key, err)
			}
			if err := s.Put(op.NewKey, val); err != nil {
				return fmt.Errorf("patch[%d] rename %q->%q put: %w", i, op.Key, op.NewKey, err)
			}
			if err := s.Delete(op.Key); err != nil {
				return fmt.Errorf("patch[%d] rename %q->%q delete: %w", i, op.Key, op.NewKey, err)
			}
		default:
			return fmt.Errorf("patch[%d]: unknown op %q", i, op.Op)
		}
	}
	return nil
}
