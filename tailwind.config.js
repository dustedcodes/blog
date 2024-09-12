const defaultTheme = require('tailwindcss/defaultTheme')

module.exports = {
  content: ["./cmd/**/*.{html,js,md}"],
  theme: {
    screens: {
      'xs': '480px',
      ...defaultTheme.screens,
    },
    fontFamily: {
      'sans': ['Inter', ...defaultTheme.fontFamily.sans],
      'serif': [...defaultTheme.fontFamily.serif],
      'mono': ['Inconsolata', ...defaultTheme.fontFamily.mono],
      'display': ['Nunito Sans', ...defaultTheme.fontFamily.sans]
    },
    colors: {
      black: '#000000',
      white: '#ffffff',
      transparent: 'transparent',
      current: 'currentColor',
      paper: '#f7f6f4',
      accent: '#cc7a3d',
      fire: '#de5553',
      ink: {
        0: '#ebe9e6',
        1: '#e3e0db',
        2: '#c5c3be',
        3: '#a8a6a2',
        4: '#8c8a87',
        5: '#71706d',
        6: '#575654',
        7: '#3f3e3c',
        8: '#272726',
        9: '#121211',
      }
    }
  },
  plugins: [],
}
