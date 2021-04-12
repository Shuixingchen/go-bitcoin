# go-bitcoin

#### 介绍
go实现简单的类比特币系统

#### 软件架构
三个主要的协程
1.挖矿-主要是挖出新区块，把交易打包到链上
2.验证交易合法性，把接收到的交易验证通过后写入待打包切片中
3.btc浏览器，与web交互

#### 安装教程
下载代码后直接运行main.go

#### 使用说明
运行程序后
本地访问btc浏览器
http://localhost:8080/blockchain //显示所有区块
http://localhost:8080/transaction/:id //显示某个区块的信息
http://localhost:8080/transaction //提交交易表单

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request
