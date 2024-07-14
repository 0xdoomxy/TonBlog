
module.exports = {
  content: ["./src/**/*.{html,js,jsx}"],
  theme: {
    extend: {},
  },
  rules: [
    {
      test: /\.scss$/,
      use: [
        'style-loader', // 将 JS 字符串生成为 style 节点
        'css-loader',   // 将 CSS 转化成 CommonJS 模块
        'sass-loader'   // 将 Sass 编译成 CSS，默认使用 Node Sass
      ]
    }
  ],
  // plugins: [
  //   require("@tailwindcss/line-clamp"),
  // ],
}

