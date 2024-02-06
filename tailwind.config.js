/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./app/templates/**/*.templ"],
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/typography")],
};
