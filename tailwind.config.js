/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./app/tmpl/**/*.templ"],
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/typography")],
};
