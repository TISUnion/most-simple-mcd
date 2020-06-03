### TODO
- 修复部分低版本服务端中文乱码
- 适配各种服务端
- godoc注释完善
- 编写golang/protobuf插件
    - 根据proto自动生成对应的gin路由
    - 根据proto自动生成接口文档
- 使用[protobuf插件](https://github.com/lightbrotherV/gin-protobuf)重写http接口
- 用golang复写sqlite, 用来替换dgraph-io/badger（危）
- 添加命令行交互模式