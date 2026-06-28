// Package feedbackimprovementloop mines sanitized recurring failures without
// mutating production behavior.
package feedbackimprovementloop

import (
	"errors"
	"sort"
)

// Correction is a sanitized failure/correction observation.
type Correction struct {
	FailureClass    string `json:"failure_class"`
	Component       string `json:"component"`
	CorrectionClass string `json:"correction_class"`
	TenantHash      string `json:"tenant_hash"`
}

// Recommendation is an aggregate proposal requiring review.
type Recommendation struct {
	FailureClass        string `json:"failure_class"`
	Component           string `json:"component"`
	CorrectionClass     string `json:"correction_class"`
	OccurrenceCount     int    `json:"occurrence_count"`
	DistinctTenants     int    `json:"distinct_tenants"`
	RequiresHumanReview bool   `json:"requires_human_review"`
}

// MineCorrections emits recommendations only when both occurrence and
// k-anonymous tenant thresholds are met.
func MineCorrections(
	corrections []Correction,
	minOccurrences, minDistinctTenants int,
) ([]Recommendation, error) {
	if minOccurrences < 2 ||
		minDistinctTenants < 2 ||
		len(corrections) > 100_000 {
		return nil, errors.New("correction mining policy is invalid")
	}
	type aggregate struct {
		count   int
		tenants map[string]struct{}
	}
	aggregates := make(map[string]*aggregate)
	records := make(map[string]Correction)
	for _, correction := range corrections {
		if correction.FailureClass == "" ||
			correction.Component == "" ||
			correction.CorrectionClass == "" ||
			correction.TenantHash == "" {
			return nil, errors.New("correction observation is incomplete")
		}
		key := correction.FailureClass + "\x00" +
			correction.Component + "\x00" +
			correction.CorrectionClass
		item := aggregates[key]
		if item == nil {
			item = &aggregate{tenants: make(map[string]struct{})}
			aggregates[key] = item
			records[key] = correction
		}
		item.count++
		item.tenants[correction.TenantHash] = struct{}{}
	}
	var recommendations []Recommendation
	for key, item := range aggregates {
		if item.count < minOccurrences || len(item.tenants) < minDistinctTenants {
			continue
		}
		record := records[key]
		recommendations = append(recommendations, Recommendation{
			FailureClass:        record.FailureClass,
			Component:           record.Component,
			CorrectionClass:     record.CorrectionClass,
			OccurrenceCount:     item.count,
			DistinctTenants:     len(item.tenants),
			RequiresHumanReview: true,
		})
	}
	sort.Slice(recommendations, func(i, j int) bool {
		if recommendations[i].OccurrenceCount != recommendations[j].OccurrenceCount {
			return recommendations[i].OccurrenceCount > recommendations[j].OccurrenceCount
		}
		return recommendations[i].FailureClass < recommendations[j].FailureClass
	})
	return recommendations, nil
}
