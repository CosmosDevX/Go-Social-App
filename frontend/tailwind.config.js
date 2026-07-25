/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        space: {
          950: '#05010d',
          900: '#0a0618',
          800: '#120b2e',
          700: '#1a1145',
          600: '#2a1a6e',
          accent: '#8b5cf6',
          cyan: '#22d3ee',
          pink: '#ec4899',
        }
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
      },
      animation: {
        'float': 'float 6s ease-in-out infinite',
        'pulse-slow': 'pulse 4s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'glow': 'glow 2s ease-in-out infinite alternate',
      },
      keyframes: {
        float: {
          '0%, 100%': { transform: 'translateY(0px)' },
          '50%': { transform: 'translateY(-12px)' },
        },
        glow: {
          '0%': { boxShadow: '0 0 20px rgba(139, 92, 246, 0.3)' },
          '100%': { boxShadow: '0 0 40px rgba(139, 92, 246, 0.6)' },
        }
      },
      backgroundImage: {
        'space-gradient': 'radial-gradient(ellipse at top, #1a1145 0%, #05010d 70%)',
        'card-gradient': 'linear-gradient(135deg, rgba(42, 26, 110, 0.4) 0%, rgba(10, 6, 24, 0.6) 100%)',
      }
    },
  },
  plugins: [],
}
