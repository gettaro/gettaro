import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, BarChart, Bar, Legend } from 'recharts'
import { GraphMetric } from '../types/memberMetrics'
import { formatTimeMetric } from '../utils/formatMetrics'

interface TeamMetricData {
  teamName: string
  metric: GraphMetric
}

interface MetricChartProps {
  metric?: GraphMetric
  teamMetrics?: TeamMetricData[]
  height?: number
}

export default function MetricChart({ metric, teamMetrics, height = 200 }: MetricChartProps) {
  // If team breakdown is enabled, combine team metrics
  let chartData: any[] = []
  let dataKeys: string[] = []
  
  if (teamMetrics && teamMetrics.length > 0) {
    // Combine multiple teams' metrics into one chart
    const allDates = new Set<string>()
    
    // Collect all dates from all teams
    teamMetrics.forEach(team => {
      team.metric.time_series?.forEach(entry => {
        allDates.add(entry.date)
      })
    })
    
    // Sort dates
    const sortedDates = Array.from(allDates).sort((a, b) => 
      new Date(a).getTime() - new Date(b).getTime()
    )
    
    // Create data points for each date with values from each team
    chartData = sortedDates.map(date => {
      const dataPoint: any = { date }
      
      teamMetrics.forEach(team => {
        const teamKey = team.teamName.replace(/\s+/g, '_').toLowerCase()
        const entry = team.metric.time_series?.find(e => e.date === date)
        if (entry && entry.data.length > 0) {
          // Use the first data point's value (assuming one value per date per team)
          dataPoint[teamKey] = entry.data[0]?.value ?? 0
        } else {
          // Set to null for missing data (Recharts will connect lines across nulls)
          dataPoint[teamKey] = null
        }
      })
      
      return dataPoint
    })
    
    // Data keys are team names
    dataKeys = teamMetrics.map(team => 
      team.teamName.replace(/\s+/g, '_').toLowerCase()
    )
  } else if (metric) {
    // Single metric (cumulative)
    chartData = metric.time_series?.map(entry => {
      const dataPoint: any = { date: entry.date }
      
      entry.data.forEach(point => {
        dataPoint[point.key] = point.value
      })
      
      return dataPoint
    }) || []

    dataKeys = Array.from(new Set(
      metric.time_series?.flatMap(entry => entry.data.map(point => point.key)) || []
    ))
  }

  // Determine chart type
  const referenceMetric = teamMetrics?.[0]?.metric || metric
  const isBarChart = referenceMetric?.type?.toLowerCase().includes('bar') || 
                     referenceMetric?.type?.toLowerCase().includes('count') ||
                     referenceMetric?.type?.toLowerCase().includes('total')

  // Check if this is a time-based metric (seconds unit)
  const isTimeMetric = referenceMetric?.unit === 'seconds' || referenceMetric?.unit === 'time'
  
  // Y-axis formatter for time metrics
  const formatYAxisTick = (value: number) => {
    if (isTimeMetric) {
      return formatTimeMetric(value)
    }
    return value.toLocaleString()
  }

  // Tooltip formatter for time metrics
  const formatTooltipValue = (value: number, name: string) => {
    if (isTimeMetric && value !== null && value !== undefined) {
      return [formatTimeMetric(value), name]
    }
    return [value.toLocaleString(), name]
  }

  const colors = [
    '#8884d8', '#82ca9d', '#ffc658', '#ff7300', '#00ff00', '#ff00ff', '#00ffff',
    '#ff6b6b', '#4ecdc4', '#45b7d1', '#f9ca24', '#f0932b', '#eb4d4b', '#6c5ce7'
  ]

  // Check if there's any actual data (not all null/undefined values)
  // Note: 0 is a valid value for metrics, so we only check for null/undefined
  const hasData = chartData.length > 0 && chartData.some(point => {
    return dataKeys.some(key => {
      const value = point[key]
      return value !== null && value !== undefined
    })
  })

  if (!hasData) {
    return null
  }

  return (
    <div className="w-full" style={{ height }}>
      <ResponsiveContainer width="100%" height="100%">
        {isBarChart ? (
          <BarChart data={chartData} margin={{ top: 5, right: 30, left: 20, bottom: 5 }}>
            <CartesianGrid strokeDasharray="3 3" stroke="hsl(var(--muted-foreground))" opacity={0.3} />
            <XAxis 
              dataKey="date" 
              stroke="hsl(var(--muted-foreground))"
              fontSize={12}
              tickFormatter={(value) => {
                const date = new Date(value)
                return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
              }}
            />
            <YAxis 
              stroke="hsl(var(--muted-foreground))" 
              fontSize={12}
              tickFormatter={formatYAxisTick}
            />
            <Tooltip 
              contentStyle={{
                backgroundColor: 'hsl(var(--background))',
                border: '1px solid hsl(var(--border))',
                borderRadius: '6px',
                color: 'hsl(var(--foreground))'
              }}
              labelFormatter={(value) => {
                const date = new Date(value)
                return date.toLocaleDateString('en-US', { 
                  year: 'numeric', 
                  month: 'short', 
                  day: 'numeric' 
                })
              }}
              formatter={formatTooltipValue}
            />
            {dataKeys.map((key, index) => (
              <Bar 
                key={key} 
                dataKey={key} 
                fill={colors[index % colors.length]}
                name={teamMetrics 
                  ? teamMetrics[index]?.teamName || key
                  : key.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())
                }
              />
            ))}
            {teamMetrics && teamMetrics.length > 0 && (
              <Legend 
                wrapperStyle={{ paddingTop: '20px' }}
                formatter={(value) => {
                  const team = teamMetrics.find(t => 
                    t.teamName.replace(/\s+/g, '_').toLowerCase() === value
                  )
                  return team?.teamName || value
                }}
              />
            )}
          </BarChart>
        ) : (
          <LineChart data={chartData} margin={{ top: 5, right: 30, left: 20, bottom: 5 }}>
            <CartesianGrid strokeDasharray="3 3" stroke="hsl(var(--muted-foreground))" opacity={0.3} />
            <XAxis 
              dataKey="date" 
              stroke="hsl(var(--muted-foreground))"
              fontSize={12}
              tickFormatter={(value) => {
                const date = new Date(value)
                return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
              }}
            />
            <YAxis 
              stroke="hsl(var(--muted-foreground))" 
              fontSize={12}
              tickFormatter={formatYAxisTick}
            />
            <Tooltip 
              contentStyle={{
                backgroundColor: 'hsl(var(--background))',
                border: '1px solid hsl(var(--border))',
                borderRadius: '6px',
                color: 'hsl(var(--foreground))'
              }}
              labelFormatter={(value) => {
                const date = new Date(value)
                return date.toLocaleDateString('en-US', { 
                  year: 'numeric', 
                  month: 'short', 
                  day: 'numeric' 
                })
              }}
              formatter={formatTooltipValue}
            />
            {dataKeys.map((key, index) => (
              <Line 
                key={key} 
                type="monotone" 
                dataKey={key} 
                stroke={colors[index % colors.length]}
                strokeWidth={2}
                dot={{ fill: colors[index % colors.length], strokeWidth: 2, r: 4 }}
                connectNulls={true}
                name={teamMetrics 
                  ? teamMetrics[index]?.teamName || key
                  : key.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())
                }
              />
            ))}
            {teamMetrics && teamMetrics.length > 0 && (
              <Legend 
                wrapperStyle={{ paddingTop: '20px' }}
                formatter={(value) => {
                  const team = teamMetrics.find(t => 
                    t.teamName.replace(/\s+/g, '_').toLowerCase() === value
                  )
                  return team?.teamName || value
                }}
              />
            )}
          </LineChart>
        )}
      </ResponsiveContainer>
    </div>
  )
}
