module.exports = {
  theme: {
    extend: {
      keyframes: {
        shake: {
          '0%, 100%': { transform: 'translateX(0)' },
          '10%, 30%, 50%, 70%, 90%': { transform: 'translateX(-5px)' },
          '20%, 40%, 60%, 80%': { transform: 'translateX(5px)' },
        },
        fadeIn: {
          '0%': { opacity: 0, transform: 'translateY(-10px)' },
          '100%': { opacity: 1, transform: 'translateY(0)' },
        },
        glow: {
          '0%, 100%': { filter: 'drop-shadow(0 0 2px #00E5B0)' },
          '50%': { filter: 'drop-shadow(0 0 8px #00E5B0) drop-shadow(0 0 12px #00B8FF)' },
        }
      },
      animation: {
        shake: 'shake 0.5s ease-in-out',
        fadeIn: 'fadeIn 0.3s ease-out',
        glow: 'glow 2s ease-in-out infinite',
      },
    },
  },
}