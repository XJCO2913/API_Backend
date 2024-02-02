## api.backend.xjco2013

SE项目后端仓库

### 项目结构

<img src="/Users/yuerfei/Documents/SWJTU/Year3/SoftwareProject/repo/backend/img/repo_tree.png" alt="项目结构" style="zoom:50%; float:left" />

#### 1. cmd

cmd是package main所在的目录，也是项目的启动目录。包含main函数，路由的初始化，以及一些项目启动时的必要资源的初始化

#### 2. config

config目录用于存放各类配置文件，如数据库配置文件(包含数据库账号密码等等)，环境变量，JWT加密密钥，docker配置文件等等。

将敏感信息以配置的形式存放，而不是直接写在代码里，提高安全性。（所以这个目录下的文件有一些不会推到仓库里）

#### 3. controller

controller层是请求到达服务器