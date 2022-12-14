const defaultTheme = require('tailwindcss/defaultTheme')

module.exports = {
  content: ["./cmd/**/*.{html,js,md}"],
  theme: {
    screens: {
      'xs': '400px',
      ...defaultTheme.screens,
    },
    fontFamily: {
      'sans': ['Inter', ...defaultTheme.fontFamily.sans],
      'serif': [...defaultTheme.fontFamily.serif],
      'mono': ['Inconsolata', ...defaultTheme.fontFamily.mono],
      'nunito': ['Nunito', ...defaultTheme.fontFamily.sans]
    },
    extend: {
    },
  },
  plugins: [],
}
