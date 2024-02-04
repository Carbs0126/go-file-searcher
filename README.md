# go-file-searcher
在mac操作系统的某个目录下，根据关键字搜索匹配该目录下的所有名称，并可以直接在命令行中打开
```
./go-file-searcher keywords-you-want-to-search
```
## 示例
![image](https://github.com/Carbs0126/go-file-searcher/assets/14228871/48ba4ca9-70bc-41c4-a16c-5735573560a7)

## 使用
1. 命令行进入要搜索的文件夹目录下；
2. 点击上下箭头，可以选中某个文件，点击``enter``键，即可使用mac操作的open命令打开该文件
3. `-help`获取使用帮助；

## 编译源码
1. 在项目根目录下使用命令``go get``；
2. 在项目根目录下使用命令``go build``，即可在项目根目录下得到 `go-file-searcher` 可执行文件；

## TODO
1. 如果终端的输出文字大于一屏幕行数，则光标向上只能定位到屏幕边缘，暂时未能实现继续向上滚屏的功能；

## 感谢
1. 开源项目
   - github.com/eiannone/keyboard
   - golang.org/x/sys
2. 参考资料
   - ANSI转义序列 https://zh.wikipedia.org/wiki/ANSI%E8%BD%AC%E4%B9%89%E5%BA%8F%E5%88%97
