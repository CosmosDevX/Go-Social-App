import { useMemo } from 'react'

export function Stars() {
  const stars = useMemo(() => {
    return Array.from({ length: 80 }, (_, i) => ({
      id: i,
      left: `${Math.random() * 100}%`,
      top: `${Math.random() * 100}%`,
      size: Math.random() * 2 + 1,
      duration: `${2 + Math.random() * 4}s`,
      delay: `${Math.random() * 3}s`,
    }))
  }, [])

  return (
    <div className="stars" aria-hidden>
      {stars.map((star) => (
        <div
          key={star.id}
          className="star"
          style={{
            left: star.left,
            top: star.top,
            width: star.size,
            height: star.size,
            // @ts-expect-error custom prop
            '--duration': star.duration,
            animationDelay: star.delay,
          }}
        />
      ))}
    </div>
  )
}
