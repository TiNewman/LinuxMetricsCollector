module.exports = {
  content: [
    './public/index.html',
    "./pages/**/*.{js,ts,jsx,tsx}",
    './src/**/*.css',
  ],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")]
}
