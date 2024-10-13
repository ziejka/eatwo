/** @type {import('tailwindcss').Config} */
export default {
  content: ["./views/**/*.templ"], // this is where our templates are located
  theme: {
    extend: {
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
