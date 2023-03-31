## webpack 中使用 alias

- 在 craco.config.js 中添加配置：

```text
  webpack: {
    alias: {
      '@': path.resolve(__dirname, 'src')
    },
    plugins: [
      // 查看打包的进度
      new SimpleProgressWebpackPlugin()
    ],

    configure: {}
  },

```

- 在 tsconfig.jso 中添加

```text
 "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"]
    }
  }
```

- 在 eslint 中配置  
  在 webpack 项目之中， 我们会借助 alias 别名提升代码效率和打包效率。
  但是在使用了自定义的路径指向后，eslint 就会对应产生找不到模块的报错。

```text
Unable to resolve path to module '@/component' import/no-unresolved
```

这个时候，就要使用 npm 包 eslint-import-resolver-alias 来解决了

```text
yarn add eslint-import-resolver-alias --dev
```

在.eslintrc.js 文件中添加

```text
// eslint配置
 "settings": {
   "import/resolver": {
     "alias": {
       "map": [
         ["@", "./src"]
       ]
     }
   }
 },
```

**注意：在 tsconfig.jso 中的生效且能被 idea 识别的 path 是合并后的 path, 这里就是"./src/\*"**

```text
"baseUrl": ".",
"paths": {
  "@/*": ["src/*"]
},
```
