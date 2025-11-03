/**
 * Formats a time value in seconds to the most appropriate unit
 * @param seconds - Time value in seconds
 * @returns Formatted time string with unit
 * 
 * @example
 * formatTimeMetric(30)     // "30s"
 * formatTimeMetric(90)     // "1.5m"
 * formatTimeMetric(3600)   // "1h"
 * formatTimeMetric(86400)  // "1d"
 * formatTimeMetric(259200) // "3d"
 * formatTimeMetric(604800) // "1w"
 * formatTimeMetric(2592000) // "1mo"
 */
export function formatTimeMetric(seconds: number): string {
  if (seconds < 60) {
    // Less than 1 minute: show as seconds
    return `${Math.round(seconds)}s`
  } else if (seconds < 3600) {
    // 1 minute to 1 hour: show as minutes
    const minutes = seconds / 60
    if (minutes < 10) {
      return `${minutes.toFixed(1)}m`
    } else if (minutes < 60) {
      return `${Math.round(minutes)}m`
    } else {
      return `${Math.round(minutes)}m`
    }
  } else if (seconds < 86400) {
    // 1 hour to 1 day: show as hours
    const hours = seconds / 3600
    if (hours < 10) {
      return `${hours.toFixed(1)}h`
    } else {
      return `${Math.round(hours)}h`
    }
  } else {
    // More than 1 day: show as days
    const days = seconds / 86400
    if (days < 10) {
      return `${days.toFixed(1)}d`
    } else if (days < 30) {
      return `${Math.round(days)}d`
    } else {
      // For very long periods, show weeks or months
      const weeks = days / 7
      if (weeks < 4) {
        return `${Math.round(weeks)}w`
      } else {
        const months = days / 30.44 // Average days per month
        return `${Math.round(months)}mo`
      }
    }
  }
}

/**
 * Formats a metric value based on its unit
 * @param value - The metric value
 * @param unit - The unit of measurement
 * @returns Formatted value string
 * 
 * @example
 * // Time metrics (automatically converted to appropriate unit)
 * formatMetricValue(30, 'seconds')     // "30s"
 * formatMetricValue(90, 'seconds')     // "1.5m"
 * formatMetricValue(3600, 'seconds')   // "1h"
 * 
 * // Count metrics
 * formatMetricValue(1234, 'count')     // "1,234"
 * 
 * // LOC metrics
 * formatMetricValue(1500, 'loc')       // "1.5k"
 * formatMetricValue(500, 'loc')        // "500"
 * 
 * // Percent metrics
 * formatMetricValue(75.5, 'percent')   // "75.5%"
 * formatMetricValue(100, 'percent')   // "100%"
 */
export function formatMetricValue(value: number, unit: string): string {
  switch (unit) {
    case 'time':
    case 'seconds':
      return formatTimeMetric(value)
    case 'count':
      return value.toLocaleString()
    case 'loc':
      if (value >= 1000) {
        return `${(value / 1000).toFixed(1)}k`
      }
      return value.toLocaleString()
    case 'percent':
      // Format percentage with appropriate decimal places
      if (value % 1 === 0) {
        return `${Math.round(value)}%`
      }
      return `${value.toFixed(1)}%`
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
    case 'seconds':
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

/**
 * Gets icon information based on the backend iconIdentifier and iconColor
 * @param iconIdentifier - The icon identifier from the backend
 * @param iconColor - The icon color from the backend
 * @returns Object with icon type and color information
 */
export function getBackendIconInfo(iconIdentifier: string, iconColor: string) {
  return {
    type: iconIdentifier,
    color: iconColor,
    colorClass: `text-${iconColor}-500`
  }
} 