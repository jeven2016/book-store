```text
declare var 声明全局变量, 还有 declare let 和 declare const
declare function 声明全局方法
declare class 声明全局类
declare enum 声明全局枚举类型
declare namespace 声明（含有子属性的）全局对象
interface 和 type 声明全局类型

```

```text
declare function jQuery(selector: string): any;


declare class Animal {
    name: string;
    constructor(name: string);
    sayHi(): string;
}

declare enum Directions {
    Up,
    Down,
    Left,
    Right
}


interface AjaxSettings {
    method?: 'GET' | 'POST'
    data?: any;
}

```

注意，只有 function、class 和 interface 可以直接默认导出，其他的变量需要先定义出来，再默认导出

发布声明文件§
当我们为一个库写好了声明文件之后，下一步就是将它发布出去了。

此时有两种方案：

- 将声明文件和源码放在一起
- 将声明文件发布到 @types 下
  这两种方案中优先选择第一种方案。保持声明文件与源码在一起，使用时就不需要额外增加单独的声明文件库的依赖了，而且也能保证声明文件的版本与源码的版本保持一致。

仅当我们在给别人的仓库添加类型声明文件，但原作者不愿意合并 pull request 时，才需要使用第二种方案，将声明文件发布到 @types 下。

**不推荐的做法，去掉强制检查**

在 compilerOptions 中添加 "noImplicitAny": false  
这里说一下 "noImplicitAny" 这个参数是将 TypeScript 从可选类型语言转换为强制类型检验语言，
简单点说就是在 JS 中有一些不怎么规范的写法都会被 TS 提示警告，所以网上大部分人其实是推荐 noImplicitAny=true 的，
这样可以规范自己的代码，提高自己的编码能力

```text
"noImplicitAny": false,
```
