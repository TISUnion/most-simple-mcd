- [/most.simple.mcd.User/login](#mostsimplemcduserlogin)  用户登录
- [/most.simple.mcd.User/logout](#mostsimplemcduserlogout)  用户注销
- [/most.simple.mcd.User/info](#mostsimplemcduserinfo)  获取用户信息
- [/most.simple.mcd.User/update](#mostsimplemcduserupdate)  更新用户信息

##/most.simple.mcd.User/login
### 用户登录

#### 方法：GRPC

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|account|否|string| 账号|
|password|否|string| 密码|


#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": 
    {
        // 登录态token
        "token": ""
    }

}
```


##/most.simple.mcd.User/logout
### 用户注销

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


##/most.simple.mcd.User/info
### 获取用户信息

#### 方法：GRPC

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|token|否|string| 登录态token|


#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": 
    {
        // 账号
        "account": ""
        // 密码
        "password": ""
        // 昵称
        "nickname": ""
        // 权限
        "roles": [
            ""
        ]
        // 头像
        "avatar": ""
    }

}
```


##/most.simple.mcd.User/update
### 更新用户信息

#### 方法：GRPC

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|account|否|string| 账号|
|password|否|string| 密码|
|nickname|否|string| 昵称|
|roles|否|string| 权限|
|avatar|否|string| 头像|


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


