package engine

import "ems.dev/backend/services/sourcecontrol/types"

// mergeTimeSeriesWithPeers merges member time series data with peer values by date
func mergeTimeSeriesWithPeers(memberSeries []types.TimeSeriesEntry, peerSeries []types.TimeSeriesEntry, memberLabel string) []types.TimeSeriesEntry {
	// Create a map of peer values by date for quick lookup
	peerMap := make(map[string]float64)
	for _, entry := range peerSeries {
		if len(entry.Data) > 0 {
			peerMap[entry.Date] = entry.Data[0].Value // Peer series always has one data point with "Peers" key
		}
	}

	// Merge peer values into member series
	merged := make([]types.TimeSeriesEntry, len(memberSeries))
	for i, entry := range memberSeries {
		mergedEntry := types.TimeSeriesEntry{
			Date: entry.Date,
			Data: make([]types.TimeSeriesDataPoint, len(entry.Data)),
		}
		
		// Copy member data points
		copy(mergedEntry.Data, entry.Data)
		
		// Add peer value if available for this date
		if peerValue, exists := peerMap[entry.Date]; exists {
			mergedEntry.Data = append(mergedEntry.Data, types.TimeSeriesDataPoint{
				Key:   "Peers",
				Value: peerValue,
			})
		}
		
		merged[i] = mergedEntry
	}

	return merged
}

