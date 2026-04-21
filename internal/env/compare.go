package env

import "fmt"

// CompareResult holds the outcome of comparing two Sets.
type CompareResult struct {
	OnlyInSrc  map[string]string // keys present only in src
	OnlyInDst  map[string]string // keys present only in dst
	Same       map[string]string // keys with identical values in both
	Conflicted map[string][2]string // keys present in both but with different values [src, dst]
}

// Compare performs a detailed comparison between src and dst Sets.
// It returns a CompareResult describing which keys are unique to each set,
// shared, or conflicting.
func Compare(src, dst *Set) (*CompareResult, error) {
	if src == nil {
		return nil, fmt.Errorf("compare: src set is nil")
	}
	if dst == nil {
		return nil, fmt.Errorf("compare: dst set is nil")
	}

	result := &CompareResult{
		OnlyInSrc:  make(map[string]string),
		OnlyInDst:  make(map[string]string),
		Same:       make(map[string]string),
		Conflicted: make(map[string][2]string),
	}

	srcVars := src.All()
	dstVars := dst.All()

	for k, sv := range srcVars {
		if dv, ok := dstVars[k]; ok {
			if sv == dv {
				result.Same[k] = sv
			} else {
				result.Conflicted[k] = [2]string{sv, dv}
			}
		} else {
			result.OnlyInSrc[k] = sv
		}
	}

	for k, dv := range dstVars {
		if _, ok := srcVars[k]; !ok {
			result.OnlyInDst[k] = dv
		}
	}

	return result, nil
}

// IsEqual returns true when src and dst contain exactly the same key-value pairs.
func IsEqual(src, dst *Set) (bool, error) {
	cr, err := Compare(src, dst)
	if err != nil {
		return false, err
	}
	return len(cr.OnlyInSrc) == 0 && len(cr.OnlyInDst) == 0 && len(cr.Conflicted) == 0, nil
}
