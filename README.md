# Hostloc 签到脚本
Hostloc 全球主机交流论坛 刷分程序 自动登录 自动访问别人空间

## 使用
1. 建一个名为`conf.yaml`的配置文件
2. 下载程序，运行。例如 `nohup ./hostloc &`
```bash
# 新建文件夹
mkdir hostloc && cd hostloc
# 下载配置文件例子 然后修改成自己的信息  下面 2选1
curl -s https://raw.githubusercontent.com/uniqueque/hostloc/main/conf.yaml.example > conf.yaml
wget -O conf.yaml https://raw.githubusercontent.com/uniqueque/hostloc/main/conf.yaml.example
# 记得修改配置文件
vim conf.yaml

#下载对应系统的运行运行，下面以linux amd64举例
wget -O hostloc.tar.gz https://github.com/uniqueque/hostloc/releases/download/v0.0.1/hostloc_0.0.1_Linux_x86_64.tar.gz
# 解压
tar -zxvf hostloc.tar.gz
# 运行
nohup ./hostloc &
```
## 配置文件

```yaml
tg_token: <TgToken的值>
tg_id: <TgId的值>
users:
  - username: <第一个用户的Username/email的值>
    password: <第一个用户的Password的值>
    fastloginfield: <登录类型[email/username]用户名不包含@的话可不填>
    tg_id: <指定别的tgid接受消息，不填就推送到外面的那个tgid>
  - username: <第二个用户的Username的值>
    password: <第二个用户的Password的值>
    fastloginfield: <第二个用户的FastLoginField的值>
    tg_id: <指定别的tgid接受消息，不填就推送到外面的那个tgid>
# 可以继续添加更多的用户

```

