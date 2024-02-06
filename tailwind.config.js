/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./tmpl/**/*.templ"],
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/typography")],
};
