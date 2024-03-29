# go-file-searcher
在mac操作系统的某个目录下，根据关键字搜索匹配该目录下的所有名称，并可以直接在命令行中打开
```
./go-file-searcher keywords-you-want-to-search  [-h] [-r] [-t]
```
`keywords-you-want-to-search`可以使用空格隔开，如 `./go-file-searcher read .md` 可以搜索当前目录下文件名（忽略大小写）匹配 `read.*\.md`正则表达式的结果；

## 示例
![image](https://github.com/Carbs0126/go-file-searcher/assets/14228871/03cff9fd-a718-4d45-a79a-5164b5009842)
![image](https://github.com/Carbs0126/go-file-searcher/assets/14228871/3a0a2faf-cf06-4f7e-afe8-ab7d68d010f7)
![image](https://github.com/Carbs0126/go-file-searcher/assets/14228871/fee612ef-8c19-440f-9488-23f4151838de)
注：上图的`se`是将`go-file-searcher`的执行文件重新命名并添加到系统路径中。进入`.gradle`文件夹下，递归遍历搜索所有文件名中带有 `gradle` 和 `all` 的文件。

## 使用
1. 命令行进入要搜索的文件夹目录下；
2. 点击上下箭头，可以选中某个文件，点击``enter``键，即可使用mac操作的open命令打开该文件，点击``space``键，即可打开该文件的目录；
3. 当文件较多，输出的行数大于一屏时，会自动切屏。此时点击左右箭头，可以切屏展示其它文件；
4. 在任何一个参数位置输入`-help`或者`-h`，将会获取使用帮助；
5. 在任何一个参数位置输入`-plain`或者`-p`，将会只搜索第一层级文件夹，否则会递归遍历当前文件夹下的所有子文件夹；
6. 在任何一个参数位置输入`-time`或者`-t`，会把结果按照时间顺序排列显示，最后更新的文件会显示在最开始位置；

## 编译源码
1. 在项目根目录下使用命令``go get``；
2. 在项目根目录下使用命令``go build``，即可在项目根目录下得到 `go-file-searcher` 可执行文件；

## TODO

1. 点击`空格`弹出菜单，进而进行选择并弹出对话框
   ```
      [F] >  android/someproject   
      [F]    +-----------------------+   +----------------------------------------+ 
      [F]    | > Info            [I] |   |              File Info                 |
      [F]    |   Rename          [R] |   |    Size:  300000 Byte (3MB)            |
      [F]    |   Delete          [D] |   |    Type:  jpeg  |  1920x1080           |
      [F]    |   Parent Folder   [P] |   |  Create:  1970.01.01 20:00:01          |
      [F]    |   Close           [C] |   +----------------------------------------+
      [F]    +-----------------------+
      [F]    abcdefg.txt
   ```
2. 对话框？
   ```
      Esc   : dismiss dialog
      Enter : select and dismiss dialog
      
      +-------------------------------------+ 
      |         Delete This File ?          |
      |                                     |
      |       [Yes]            No           |
      +-------------------------------------+
      
      +----------------------------------------+ 
      |              File Info                 |
      |    Size:  300000 Byte (3MB)            |
      |    Type:  jpeg  |  1920x1080           |
      |  Create:  1970.01.01 20:00:01          |
      +----------------------------------------+
   ```
## 感谢
1. 开源项目
   - github.com/eiannone/keyboard
   - golang.org/x/sys
2. 参考资料
   - ANSI转义序列 https://zh.wikipedia.org/wiki/ANSI%E8%BD%AC%E4%B9%89%E5%BA%8F%E5%88%97
