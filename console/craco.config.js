const { when, whenDev, whenProd, whenTest, ESLINT_MODES, POSTCSS_MODES } = require('@craco/craco');

const path = require('path');

const SimpleProgressWebpackPlugin = require('simple-progress-webpack-plugin');
const CircularDependencyPlugin = require('circular-dependency-plugin');
const { BundleAnalyzerPlugin } = require('webpack-bundle-analyzer');

module.exports = {
  webpack: {
    alias: {
      '@': path.resolve(__dirname, 'src')
    },
    plugins: [
      // 查看打包的进度
      new SimpleProgressWebpackPlugin(),

      // 新增模块循环依赖检测插件
      ...whenDev(
        () => [
          new CircularDependencyPlugin({
            exclude: /node_modules/,
            include: /src/,
            failOnError: true,
            allowAsyncCycles: false,
            cwd: process.cwd()
          })
        ],
        []
      ),
      // 新增打包分析插件
      ...whenProd(
        () => [
          // https://www.npmjs.com/package/webpack-bundle-analyzer
          new BundleAnalyzerPlugin({
            analyzerMode: 'static', // html 文件方式输出编译分析
            openAnalyzer: false,
            reportFilename: path.resolve(__dirname, `analyzer/index.html`)
          })
        ],
        []
      )
    ],

    configure: {}
  },
  jest: {
    configure: {
      moduleNameMapper: {
        '@': path.resolve(__dirname, 'src')
      }
    }
  },
  //配置接口跨域代理
  devServer: {
    port: 8088,
    proxy: {
      // '/api': 'http://localhost:9001'
      '/api': 'http://localhost:8080'
      // '/api': {
      //   target: 'http://localhost:9999',
      //   changeOrigin: true,
      //   pathRewrite: {
      //     '^/api': '/api'
      //   }
      // }
    }
  }
};
