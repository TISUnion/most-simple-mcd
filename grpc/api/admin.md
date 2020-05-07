- [/most.simple.mcd.Admin/getConfig](#mostsimplemcdadmingetconfig)  获取配置
- [/most.simple.mcd.Admin/updateConfig](#mostsimplemcdadminupdateconfig)  修改配置
- [/most.simple.mcd.Admin/operatePlugin](#mostsimplemcdadminoperateplugin)  服务端插件操作
- [/most.simple.mcd.Admin/getConfigVal](#mostsimplemcdadmingetconfigval)  服务端插件操作
- [/most.simple.mcd.Admin/runCommand](#mostsimplemcdadminruncommand)  向服务端执行一条命令
- [/most.simple.mcd.Admin/getLog](#mostsimplemcdadmingetlog)  获取服务端日志
- [/most.simple.mcd.Admin/delTmpFlie](#mostsimplemcdadmindeltmpflie)  删除临时文件
- [/most.simple.mcd.Admin/addUpToContainer](#mostsimplemcdadminadduptocontainer)  获取上传服务端文件，并注入到容器中
- [/most.simple.mcd.Admin/closeMcd](#mostsimplemcdadminclosemcd) 

##/most.simple.mcd.Admin/getConfig
### 获取配置

#### 方法：GRPC

#### 请求参数



#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": 
    {
        //
        "list": [
            {
                // 配置名称
                "config_val": ""
                // 配置值
                "config_key": ""
                // 配置等级
                "level": 0
                // 配置描述
                "description": ""
                // 是否在后台可修改
                "is_alterable": false
            }
        ]
    }

}
```


##/most.simple.mcd.Admin/updateConfig
### 修改配置

#### 方法：GRPC

#### 请求参数

```javascript
    {
        //
        "list": [
            {
                // 配置名称
                "config_val": ""
                // 配置值
                "config_key": ""
                // 配置等级
                "level": 0
                // 配置描述
                "description": ""
                // 是否在后台可修改
                "is_alterable": false
            }
        ]
    }

```



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


##/most.simple.mcd.Admin/operatePlugin
### 服务端插件操作

#### 方法：GRPC

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|server_id|否|string||
|plugin_id|否|string||
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


##/most.simple.mcd.Admin/getConfigVal
### 服务端插件操作

#### 方法：GRPC

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|name|否|string| 配置名称|


#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": 
    {
        // 配置名称
        "config_val": ""
        // 配置值
        "config_key": ""
        // 配置等级
        "level": 0
        // 配置描述
        "description": ""
        // 是否在后台可修改
        "is_alterable": false
    }

}
```


##/most.simple.mcd.Admin/runCommand
### 向服务端执行一条命令

#### 方法：GRPC

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|command|否|string| 命令|
|id|否|string||
|type|否|integer| 1：插件运行命令  2：服务端运行命令   3：插件、服务端都运行|


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


##/most.simple.mcd.Admin/getLog
### 获取服务端日志

#### 方法：GRPC

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|type|否|integer| 日志类型 1. 根据id获取服务端日志 2. gin服务器日志 3. 默认全局日志|
|id|否|string| 如果根据id获取服务端日志，则需要传id|


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


##/most.simple.mcd.Admin/delTmpFlie
### 删除临时文件

#### 方法：GRPC

#### 请求参数



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


##/most.simple.mcd.Admin/addUpToContainer
### 获取上传服务端文件，并注入到容器中

#### 方法：GRPC

#### 请求参数



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


##/most.simple.mcd.Admin/closeMcd

#### 方法：GRPC

#### 请求参数



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


