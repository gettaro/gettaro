/**
 * Formats a time value in seconds to the most appropriate unit
 * @param seconds - Time value in seconds
 * @returns Formatted time string with unit
 */
export function formatTimeMetric(seconds: number): string {
  if (seconds < 60) {
    return `${Math.round(seconds)}s`
  } else if (seconds < 3600) {
    const minutes = seconds / 60
    return `${minutes < 10 ? minutes.toFixed(1) : Math.round(minutes)}m`
  } else if (seconds < 86400) {
    const hours = seconds / 3600
    return `${hours < 10 ? hours.toFixed(1) : Math.round(hours)}h`
  } else {
    const days = seconds / 86400
    return `${days < 10 ? days.toFixed(1) : Math.round(days)}d`
  }
}

/**
 * Formats a metric value based on its unit
 * @param value - The metric value
 * @param unit - The unit of measurement
 * @returns Formatted value string
 */
export function formatMetricValue(value: number, unit: string): string {
  switch (unit) {
    case 'time':
      return formatTimeMetric(value)
    case 'count':
      return value.toLocaleString()
    case 'loc':
      if (value >= 1000) {
        return `${(value / 1000).toFixed(1)}k`
      }
      return value.toLocaleString()
    default:
      return value.toLocaleString()
  }
}

/**
 * Gets an appropriate icon for a metric based on its unit and label
 * @param unit - The unit of measurement
 * @param label - The metric label
 * @returns SVG icon component or icon name
 */
export function getMetricIcon(unit: string, label: string): string {
  switch (unit) {
    case 'time':
      return '⏱️'
    case 'count':
      if (label.toLowerCase().includes('pr')) {
        return '🔀'
      } else if (label.toLowerCase().includes('review')) {
        return '👁️'
      } else if (label.toLowerCase().includes('comment')) {
        return '💬'
      }
      return '📊'
    case 'loc':
      return '📝'
    default:
      return '📈'
  }
} 