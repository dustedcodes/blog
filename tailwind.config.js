const defaultTheme = require('tailwindcss/defaultTheme')

module.exports = {
  content: ["./cmd/**/*.{html,js,md}"],
  theme: {
    fontFamily: {
      'sans': [...defaultTheme.fontFamily.sans],
      'serif': [...defaultTheme.fontFamily.serif],
      'mono': [...defaultTheme.fontFamily.mono],
      'nunito': ['Nunito', ...defaultTheme.fontFamily.sans]
  },
    extend: {},
  },
  plugins: [],
}
