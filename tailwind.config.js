/** @type {import('tailwindcss').Config} */
export default {
  content: ["./views/**/*.templ"], // this is where our templates are located
  theme: {
    extend: {
      colors: {
        mirage: {
          '50': '#f2f5fc',
          '100': '#e3e9f6',
          '200': '#cdd9f0',
          '300': '#aabfe6',
          '400': '#829ed8',
          '500': '#6480cd',
          '600': '#5066c0',
          '700': '#4655af',
          '800': '#3e488f',
          '900': '#363e72',
          '950': '#242846',
        },
        'the-blue': {
          '50': '#effaff',
          '100': '#dbf3fe',
          '200': '#bfebfe',
          '300': '#92e0fe',
          '400': '#5fccfb',
          '500': '#3ab1f7',
          '600': '#2997ed',
          '700': '#1c7cd9',
          '800': '#1d63b0',
          '900': '#1d558b',
          '950': '#173454',
        },
        'wild-sand': {
          '50': '#f8f8f8',
          '100': '#f3f3f3',
          '200': '#dcdcdc',
          '300': '#bdbdbd',
          '400': '#989898',
          '500': '#7c7c7c',
          '600': '#656565',
          '700': '#525252',
          '800': '#464646',
          '900': '#3d3d3d',
          '950': '#292929',
        },
        'storm-gray': {
          '50': '#f7f7f8',
          '100': '#eeeef0',
          '200': '#d8d8df',
          '300': '#b7b6c3',
          '400': '#8f8fa1',
          '500': '#727186',
          '600': '#5c5b6e',
          '700': '#4c4a5a',
          '800': '#41404c',
          '900': '#393842',
          '950': '#26252c',
        },
      },
      fontFamily: {
        serif: ["Junge", "serif"],
      },
      boxShadow: {
        "mobile-nav": "0px -4px 8px 0px rgba(0,0,0,0.25)",
      }
    },
  },
  plugins: [],
}
