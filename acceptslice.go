package emissione

import (
	"strconv"
	"strings"
)

// AcceptSlice attaches the methods of Interface to []string, sorting in increasing order by its quality value.
type AcceptSlice []string

func (p AcceptSlice) Len() int           { return len(p) }
func (p AcceptSlice) Less(i, j int) bool { return q(p[i]) < q(p[j]) }
func (p AcceptSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func q(val string) float32 {
	// If it does not contain quality, assume 1
	if !strings.Contains(val, "q=") {
		return 1
	}

	splits := strings.Split(val, "q=")

	qString := splits[len(splits)-1]

	q, err := strconv.ParseFloat(qString, 32)
	if err != nil {
		// not parsable - be a bit lenient, and assume 0.5
		return 0.5
	}

	return float32(q)
}
