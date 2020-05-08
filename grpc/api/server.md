- [/most.simple.mcd.McServer/list](#mostsimplemcdmcserverlist)  获取服务端信息列表
- [/most.simple.mcd.McServer/getServerState](#mostsimplemcdmcservergetserverstate)  获取服务端信息列表
- [/most.simple.mcd.McServer/detail](#mostsimplemcdmcserverdetail)  获取服务端详情
- [/most.simple.mcd.McServer/operateServer](#mostsimplemcdmcserveroperateserver)  操作服务端
- [/most.simple.mcd.McServer/updateServerInfo](#mostsimplemcdmcserverupdateserverinfo)  修改服务端参数

##/most.simple.mcd.McServer/list
### 获取服务端信息列表

#### 方法：GRPC

#### 请求参数



#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": 
    {
        // 服务端配置
        "list": [
            {
                // 服务端实例id
                "id": ""
                // 服务器名称
                "name": ""
                // 执行的完整命令
                "cmd_str": [
                    ""
                ]
                // 启动服务器端口
                "port": 0
                // 运行所在工作区间
                "run_rath": ""
                // 是否是镜像服务器
                "is_mirror": false
                // 是否启动资源监听器
                "is_start_monitor": false
                // 使用内存大小，单位M
                "memory": 0
                // 服务端版本
                "version": ""
                // 服务器模式
                "game_type": ""
                // 启动状态：0.未启动 1.启动  -1.正在启动 -2.正在关闭
                "state": 0
                // 本机的ip
                "ips": [
                    ""
                ]
            }
        ]
    }

}
```


##/most.simple.mcd.McServer/getServerState
### 获取服务端信息列表

#### 方法：GRPC

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|否|string| 服务端id|


#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": 
    {
        // 启动状态：0.未启动 1.启动  -1.正在启动 -2.正在关闭
        "state": 0
    }

}
```


##/most.simple.mcd.McServer/detail
### 获取服务端详情

#### 方法：GRPC

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|否|string| 服务端id|


#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": 
    {
        // 服务端实例id
        "id": ""
        // 服务器名称
        "name": ""
        // 执行的完整命令
        "cmd_str": [
            ""
        ]
        // 启动服务器端口
        "port": 0
        // 运行所在工作区间
        "run_rath": ""
        // 是否是镜像服务器
        "is_mirror": false
        // 是否启动资源监听器
        "is_start_monitor": false
        // 使用内存大小，单位M
        "memory": 0
        // 服务端版本
        "version": ""
        // 服务器模式
        "game_type": ""
        // 启动状态：0.未启动 1.启动  -1.正在启动 -2.正在关闭
        "state": 0
        // 本机的ip
        "ips": [
            ""
        ]
        // 插件信息
        "plugin_info": [
            {
                // 插件名称
                "name": ""
                // 插件id
                "id": ""
                // 是否被禁用
                "is_ban": false
                // 使用命令
                "command_name": ""
                // 描述
                "description": ""
                // 使用帮助
                "help_description": ""
            }
        ]
    }

}
```


##/most.simple.mcd.McServer/operateServer
### 操作服务端

#### 方法：GRPC

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|否|string||
|operate_type|否|integer||


#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": 
    {
    }

}
```


##/most.simple.mcd.McServer/updateServerInfo
### 修改服务端参数

#### 方法：GRPC

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|否|string||
|name|否|string| 服务器名称|
|cmd_str|否|string||
|port|否|integer| 启动服务器端口|
|run_rath|否|string||
|is_mirror|否|bool||
|is_start_monitor|否|bool||
|memory|否|integer| 使用内存大小，单位M|
|version|否|string| 服务端版本|
|game_type|否|string||
|state|否|integer| 启动状态：0.未启动 1.启动  -1.正在启动 -2.正在关闭|
|ips|否|string| 本机的ip|


#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": 
    {
    }

}
```


