package api

import (
	"testing"

	"ems.dev/backend/services/directs/types"
	"github.com/stretchr/testify/assert"
)

// Helper functions for tests

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func verifyOrgChartNode(t *testing.T, expected, actual []types.OrgChartNode) {
	assert.Equal(t, len(expected), len(actual))
	for i := range expected {
		assert.Equal(t, expected[i].Member.ID, actual[i].Member.ID)
		assert.Equal(t, expected[i].Depth, actual[i].Depth)
		verifyOrgChartNode(t, expected[i].DirectReports, actual[i].DirectReports)
	}
}
