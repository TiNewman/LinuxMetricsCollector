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
      themes:[
        {
        mytheme: {
         "primary": "#0D5090",
         "secondary": "#FFA207",
         "accent": "#BFAFFF",
         "neutral": "#E5E7EB",
         "base-100": "#111827",
         "info": "#67E8F9",
         "success": "#57D03E",
         "warning": "#FFE202",
         "error": "#B91705",
          }
        }
      ],
      base: true,
      utils: true,
      logs: true,
      rtl: false,
      prefix: "",
      darkTheme: "dark"
    }
}
