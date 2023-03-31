### 创建工程

- create-react-app 文档地址：https://create-react-app.dev/docs/production-build
- 创建 typescript 工程

```javascript
原文链接：https://blog.csdn.net/u011521203/article/details/87546146
yarn create react-app my-app --template typescript
```

- init eslint 并得到.eslintrc.json 文件

```javascript
yarn run eslint --init
```

- 修改 eslint 配置文件
- 安装 prettier 格式化插件

```javascript
yarn add eslint-config-prettier eslint-plugin-prettier prettier --dev
```

**eslint-config-prettier** - Turns off all rules that are unnecessary or might conflict with Prettier.  
**eslint-plugin-prettier** - Runs Prettier as an ESLint rule  
更新 eslint 配置文件，添加 "plugin:prettier/recommended"

- 创建 prettier 文件，并添加配置

```javascript
{
  "semi": true,
  "tabWidth": 2,
  "printWidth": 100,
  "singleQuote": true,
  "trailingComma": "none",
  "jsxBracketSameLine": true
}
```

在 package.json 中添加命令

```javascript
"lint": "eslint ./src",
 "lint:fix": "eslint --fix ./src",
 "format": "prettier --write './**/*.{js,jsx,ts,tsx,css,md,json}' --config ./.prettierrc"
```

### 在工程中使用环境变量

```javascript
.env: Default.
.env.local: Local overrides. This file is loaded for all environments except test.
.env.development, .env.test, .env.production: Environment-specific settings.
.env.development.local, .env.test.local, .env.production.local: Local overrides of environment-specific settings.
```

Note: create-react-app 创建的工程变量名必须以 REACT*APP 开头，单词大写，以 *下划线分割，否则会无法读取。

- 变量文件的加载的优先级

```javascript
npm start: .env.development.local, .env.local, .env.development, .env
npm run build: .env.production.local, .env.local, .env.production, .env
npm test: .env.test.local, .env.test, .env (note .env.local is missing)
```

- 然后可以在变量文件中添加变量定义

```javascript
DOMAIN = www.example.com;
REACT_APP_FOO = $DOMAIN / foo;
REACT_APP_BAR = $DOMAIN / bar;
```

- 运行时指定模式

```javascript
//使用变量而非变量文件时
REACT_APP_NOT_SECRET_CODE=abcdef npm start

//使用变量文件
NODE_ENV=production npm start
```

- 代码中可以使用

```javascript
if (process.env.NODE_ENV !== 'production') {
  analytics.disable();
}
```

- html 中可以使用

```javascript
<title>%REACT_APP_WEBSITE_NAME%</title>
```

### 使用 craco 去覆盖 create-react-app 的 webpack 配置

refer to : https://github.com/gsoft-inc/craco/blob/master/packages/craco/README.md#cra-toolchain-for-beginners
使用 CRA 脚手架创建的项目，如果想要修改编译配置，通常可能会选择 npm run eject 弹出配置后魔改。但是，eject
是不可逆操作，弹出配置后，你将无法跟随官方的脚步去升级项目的 react-script 版本。
如果想要无 eject 重写 CRA 配置，目前成熟的是下面这几种方式

- 通过 CRA 官方支持的 --scripts-version 参数，创建项目时使用自己重写过的 react-scripts 包
- 使用 react-app-rewired + customize-cra 组合覆盖配置
- 使用 craco 覆盖配置

```text
 //idea多行编辑
 ALT + shift + insert 进入多行编辑模式
```

### Others

- create-react-app 已经支持热更新，不需要自行设置

### 首页图片

下载地址：https://pixabay.com/photos/woman-floating-portrait-attractive-1866858/

https://pixabay.com/photos/asia-woman-washing-clothes-stream-1807538/

## TODO

- Redux Thunk Redux 的异步处理中间件
- Redux Saga Redux 中间件,用于管理应用程序 Side Effect(副作用，例如异步获取数据，访问浏览器缓存等)
- react-copy-to-clipboard 基于 React 的复制到剪切板组件
- qrcode.react 基于 React 的生成二维码的组件
- react-syntax-highlighter 基于 React 的代码高亮组件
- react-highlight-words 基于 React 的关键字高亮

AntV 包含 G2、G6、F2、L7 以及一套完整的图表使用和设计规范, 提供强大的数据可视化需求

G2Plot 基于 G2 封装的开箱即用的可视化组件库

recharts 使用 React 和 D3 构建的自定义的图表库

Viser 支持多种主流框架的可视化库
react-image-crop 强大的图片裁切库

react-sparklines 基于数据自动生成趋势线

dom-to-image 基于 dom 生成图片的 canvas 库

react-img-editor 图片编辑器

braft-editor 富文本编辑器

powerNice markdown/富文本编辑器

GGEditor 可视化图编辑器

react-codemirror2 代码编辑器

jsoneditor json 编辑器

h5-dooring H5 页面编辑器
react-amap 高德地图插件

@uiw/react-baidu-map 百度地图
