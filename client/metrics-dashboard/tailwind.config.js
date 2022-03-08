module.exports = {
  content: [
    './public/index.html',
    "./pages/**/*.{js,ts,jsx,tsx}",
    './src/**/*.css',
    "./components/**/*.{js,ts,jsx,tsx}"
  ],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
  daisyui: {
      styled: true,
      themes: ["luxury"],
      base: true,
      utils: true,
      logs: true,
      rtl: false,
      prefix: "",
      darkTheme: "dark"
    }
}
