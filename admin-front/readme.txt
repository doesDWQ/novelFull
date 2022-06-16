开发工具：webstorm

初始化参考：https://www.tslang.cn/docs/handbook/react-&-webpack.html

安装webpack
npm install --save-dev webpack-cli
"scripts": {
  "build": "webpack --config webpack.config.js"
}

添加依赖
npm install --save react react-dom @types/react @types/react-dom

开发时依赖
npm install --save-dev typescript awesome-typescript-loader source-map-loader
awesome-typescript-loader：webpack编译typescript插件
source-map-loader： source-map-loader使用TypeScript输出的sourcemap文件来告诉webpack何时生成 自己的sourcemaps