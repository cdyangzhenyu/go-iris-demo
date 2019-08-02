# GO Iris Demo
###### 采用go语言主流web框架iris编写，实现MVC+前后端分离架构

#### 项目介绍
- 采用 iris 框架提供web后台框架
- 采用了 gorm 数据库模块 和 jwt 的单点登陆认证方式
- 默认使用 sqlite3 数据库
---

#### 项目目录结构
- apidoc 接口文档目录
- config 项目配置文件目录
- controllers 控制器文件目录
- database 数据库文件目录
- middleware 中间件文件目录
- models 模型文件目录
- tmp 测试数据库 sqlite3 文件目录
- tools 其他公用方法目录
---

#### api项目初始化

>拉取项目

```
```

>加载依赖管理包

```
export GOPROXY=https://goproxy.io
export GO111MODULE=on
cd fpga-bms-server
go get -d -v ./...
```

>项目配置文件

```
mkdir /etc/bms
cp config/config.ini /etc/bms/
```

>运行项目 

```
go run main.go // go 命令

或者：
go build main.go
./main
```

>发布

```
go build main.go
./main
```

---
##### 接口访问举例
>鉴权，默认接口访问用户名system，密码P@ssw0rd!

```
curl -i -X POST http://localhost:8080/v1/admin/login -H "Content-type: application/json" -d '{"username": "system","password": "P@ssw0rd!"}'
HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Vary: Origin
Date: Thu, 01 Aug 2019 15:18:57 GMT
Content-Length: 191

{"status":true,"msg":"登陆成功","data":{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjQ2NzYzMzcsImlhdCI6MTU2NDY3MjczN30.liYZoJEGAsE8HzeaIU27cPD_CmGJqbNiRvY3_CXBPjo"}}
```

>访问

```
curl -i -X GET http://localhost:8080/v1/admin/users -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjQ2NzYzMzcsImlhdCI6MTU2NDY3MjczN30.liYZoJEGAsE8HzeaIU27cPD_CmGJqbNiRvY3_CXBPjo"
HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Vary: Origin
Date: Thu, 01 Aug 2019 15:21:53 GMT
Content-Length: 503

{"status":true,"msg":"操作成功","data":[{"ID":2,"CreatedAt":"2019-08-01T22:37:10.919005+08:00","UpdatedAt":"2019-08-01T22:37:10.919005+08:00","DeletedAt":null,"Name":"admin","Username":"admin","Password":"$2a$10$ZRS.FSHyP.phi1Fg.NSLHeVKCYBd9HmfZ0j0cZFVXZs5rMDroKjjG","RoleID":1,"Role":{"ID":1,"CreatedAt":"2019-08-01T22:37:10.751051+08:00","UpdatedAt":"2019-08-01T22:37:10.751051+08:00","DeletedAt":null,"Name":"admin","DisplayName":"超级管理员","Description":"超级管理员","Perms":null}}]}
```

---
##### 单元测试 
>http test

```
 cd test
 go test -v  //所有测试
 
 go test -run TestUserCreate -v //单个测试
```

---

##### api 文档使用
自动生成文档 (访问过接口就会自动成功)
因为原生的 jquery.min.js 里面的 cdn 是使用国外的，访问很慢。
有条件的可以开个 vpn ,如果没有可以根据下面的方法修改一下，访问就很快了
>打开 apidoc/index.html 修改里面的

```
https://ajax.googleapis.com/ajax/libs/jquery/2.1.3/jquery.min.js

国内的 cdn

https://cdn.bootcss.com/jquery/2.1.3/jquery.min.js
```

>访问文档，从浏览器直接打开 apidoc/index.html 文件

---


#### 前端初始化

>拉取项目

```
```

>安装依赖

```
npm install
```

> 在src/utils/ 下面新建文件 apiUrl.js

> 复制内容到文件内

```
const api_url = 'http://localhost:8080'
export default api_url

```

>启动项目

```
npm run dev
```

#### 登录项目
输入地址 http://localhost


项目管理员账号 ： admin

项目管理员密码 ： admin
