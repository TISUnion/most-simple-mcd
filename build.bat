chcp 65001

:: 更新子模块
git submodule init
git submodule update
echo ---------------下载git子模块成功-----------------

:: 安装前端依赖
call npm install web-admin -S -g
echo ---------------安装前端依赖成功-------------------

:: 编译前端模块
call npm -prefix web-admin run  build:prod
echo ---------------打包前端模块成功-------------------

:: 安装go mod 包
go mod tidy
echo ---------------安装go mod 依赖成功----------------

::  将静态文件打包到可执行文件中
go build pack-webfile/compress-file.go
compress-file
echo ---------------打包前端静态文件成功----------------

:: 编译运行
go build main.go
echo ---------------编译运行文件成功----------------