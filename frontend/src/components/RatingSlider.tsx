import React from 'react'

interface RatingSliderProps {
  id: string
  value: number
  onChange: (value: number) => void
  min?: number
  max?: number
  step?: number
  disabled?: boolean
  className?: string
}

export default function RatingSlider({ 
  id, 
  value, 
  onChange, 
  min = 1, 
  max = 5, 
  step = 1, 
  disabled = false,
  className = ""
}: RatingSliderProps) {
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    onChange(parseInt(e.target.value))
  }

  const getRatingLabel = (rating: number) => {
    switch (rating) {
      case 1: return 'Poor'
      case 2: return 'Fair'
      case 3: return 'Good'
      case 4: return 'Very Good'
      case 5: return 'Excellent'
      default: return `${rating}/5`
    }
  }

  return (
    <div className={`space-y-2 ${className}`}>
      <div className="flex items-center justify-between">
        <span className="text-sm font-medium text-foreground">
          {getRatingLabel(value)}
        </span>
        <span className="text-sm text-muted-foreground">
          {value}/{max}
        </span>
      </div>
      
      <div className="relative">
        <input
          type="range"
          id={id}
          min={min}
          max={max}
          step={step}
          value={value}
          onChange={handleChange}
          disabled={disabled}
          className="w-full h-2 bg-muted rounded-lg appearance-none cursor-pointer slider"
          style={{
            background: `linear-gradient(to right, hsl(var(--primary)) 0%, hsl(var(--primary)) ${((value - min) / (max - min)) * 100}%, hsl(var(--muted)) ${((value - min) / (max - min)) * 100}%, hsl(var(--muted)) 100%)`
          }}
        />
        
        {/* Rating dots */}
        <div className="flex justify-between mt-1">
          {Array.from({ length: max - min + 1 }, (_, i) => {
            const rating = min + i
            return (
              <div
                key={rating}
                className={`w-2 h-2 rounded-full transition-colors ${
                  rating <= value 
                    ? 'bg-primary' 
                    : 'bg-muted'
                }`}
              />
            )
          })}
        </div>
      </div>
    </div>
  )
}
