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
      'display': ['Nunito', ...defaultTheme.fontFamily.sans]
    },
    colors: {
      transparent: 'transparent',
      current: 'currentColor',
      white: '#ffffff',
      'plum': {
        100: '#efedf2',
        200: '#dfdce5',
        300: '#cfcbd8',
        400: '#bfbacb',
        500: '#b0a9be',
        600: '#a199b2',
        700: '#9289a5',
        800: '#837999',
        900: '#75698c',
        1000: '#675a80',
        1100: '#594d6e',
        1200: '#4b415d',
        1300: '#3d354d',
        1400: '#30293d',
        1500: '#231e2e',
        1600: '#17131f',
        1700: '#0c0912'
      },
      'peach': {
        50: '#ffe7e6',
        100: '#ffdede',
        200: '#ffd6d5',
        300: '#ffcecd',
        400: '#ffc6c5',
        500: '#ffbebd',
        600: '#ffb5b5',
        700: '#ffadad',
        800: '#de9696',
        900: '#be7f7f',
      },
      'goldilocks': '#ffca3a',
      'shimmer': '#ffd570',
      'mint': '#caffbf',
    },
    // extend: {
    //   colors: {
    //     transparent: 'transparent',
    //     current: 'currentColor',
    //     'sinister': {
    //       50: '#e4e4e3',
    //       100: '#c9cac8',
    //       200: '#afb1ad',
    //       300: '#959893',
    //       400: '#7d7f7a',
    //       500: '#656862',
    //       600: '#4e524b',
    //       700: '#383c35',
    //       800: '#242820',
    //       900: '#11150d'
    //     }
    //   },
    // },
  },
  plugins: [],
}
