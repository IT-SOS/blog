#!/bin/bash

# 构建与部署

if [ ! -x "/usr/bin/expect" ]; then
  echo "The executable command expect is required."
  exit 1
fi

basedir=$(dirname $0)
config=$basedir/config

is_exit=0

# configrc 配置
if [ ! -e "$config/configrc" ]; then
  is_exit=1
  mkdir -p $config
  echo "<<<<<<<<<<<<<<"
  echo "Please modify the configuration file and execute again: $config/configrc"
  echo ">>>>>>>>>>>>>>"
  cat > $config/configrc <<EOF
## 常规配置
# 远程服务器地址
HOST=host
# 远程服务器ssh端口
PORT=port
# 远程服务器用户
USER=user
# 远程服务器密码
PASSWORD=password

## 如果在 /ssh/ssh/ssh_config 进行过配置且已经公钥授权（无密码登录），此项适用于通过跳板机远程部署
# 如果此处配置了，将忽略远程服务器地址、端口等配置
JUMP_HOST_NAME=jump_host_name

# 项目名
PROJECT_NAME=studynotes
# 项目远程目录
GO_PROJ_DIST=/data1/supervisord/study-notes/go
# 静态资源远程目录
STATIC_PROJ_DIST=/data1/supervisord/study-notes/html
EOF
fi

# config.yaml 配置
if [ ! -e "$config/config.yaml" ]; then
  is_exit=1
  echo "<<<<<<<<<<<<<<"
  echo "Please modify the configuration file and execute again: $config/config.yaml"
  echo ">>>>>>>>>>>>>>"
  cp $basedir/../config.yaml $config/config.yaml
fi

options="deploy_full | deploy_go_and_vue | deploy_linux_go | deploy_vue | deploy_supervisor | deploy_views | deploy_config | restart_supervisor | restart_nginx | build_darwin | build_windows"
if [ $is_exit -eq 1 ]; then
  echo "Usage: $0 $options"
  exit 1
fi

# 加载配置
. $config/configrc

# supervisord.conf 配置
update_supervisor_config() {
  tee $config/supervisord.conf <<EOF
[program:$PROJECT_NAME]
command=$GO_PROJ_DIST/bin/$PROJECT_NAME -w $GO_PROJ_DIST
user=$USER
autostart=true
autorestart=true
redirect_stderr=true
stdout_logfile=$GO_PROJ_DIST/logs/$PROJECT_NAME.log
stopsignal=15
stopasgroup=true
killasgroup=true
EOF
}

# shellcheck disable=SC1073
AskPassword=$(cat <<-EOF
expect {
  "*请输入密码*" { send "$PASSWORD\r" }
  "*password*" { send "$PASSWORD\r" }
  "#|$" { send "echo Hello.\r" }
}
EOF
)

# 执行远程命令
remoteShell() {
  /usr/bin/expect <<-EOF
set timeout -1
if { "$JUMP_HOST_NAME" == "" } {
  spawn ssh $2 -p$PORT $USER@$HOST "$1"
} else {
  spawn ssh $2 $JUMP_HOST_NAME "$1"
}

$AskPassword
$AskPassword

expect eof
EOF
}

# 上传文件
remoteSftp() {
  /usr/bin/expect <<-EOF
set timeout -1
if { "$JUMP_HOST_NAME" == "" } {
  spawn sftp -p$PORT $USER@$HOST
} else {
  spawn sftp $JUMP_HOST_NAME
}

$AskPassword

expect "sftp>"
send "put $1 $2\r"
expect "sftp>"
send "bye\r"
expect eof
EOF
}

# 验证预期关键字是否存在
checkExists() {
  echo $1|grep -ow "{succ}"|wc -l
}

PROJECT=$basedir/run/$PROJECT_NAME

build() {
  mkdir -p $basedir/run
  CGO_ENABLED=0 GOOS=$1 go build -v -o $PROJECT $basedir/../cmd/$PROJECT_NAME"_sock/"
}

run() {
  case $1 in
  "deploy_full")
    # 全量部署包含 supervisor / go / vue 服务与应用
    run deploy_supervisor
    run deploy_go_and_vue
  ;;
  "deploy_go_and_vue")
    # 部署 go / vue 资源，并重启相关服务
    run deploy_linux_go
    run deploy_views
    run deploy_config
    run restart_supervisor
    run deploy_vue
    run restart_nginx
  ;;
  "deploy_views")
    # 上送静态模板
    remoteShell "mkdir -p $GO_PROJ_DIST/web && rm -rf $GO_PROJ_DIST/web/*"
    remoteSftp "-r $basedir/../web/views" $GO_PROJ_DIST/web
  ;;
  "deploy_supervisor")
    echo "deploy supervisor start"
    # 检查 supervisor 是否安装
    ok=$(remoteShell "which supervisorctl && echo {succ}")
    ok=$(echo $ok|grep -ow "{succ}"|wc -l)
    if [ "$ok" -lt "2" ]; then
      echo "Not found supervisorctl. Please apt-get install supervisor."
      exit 1
    fi
    # 生成配置文件
    update_supervisor_config
    # 上送配置
    remoteSftp $config/supervisord.conf /tmp/supervisord.conf
    # 配置文件归位
    remoteShell "sudo mv /tmp/supervisord.conf /etc/supervisor/conf.d/$PROJECT_NAME.conf" "-t"
    echo "deploy supervisor end"
  ;;
  "restart_supervisor")
    # 重启 supervisor
    remoteShell "sudo systemctl restart supervisor" "-t"
  ;;
  "restart_nginx")
    # 重启 nginx
    remoteShell "sudo systemctl restart nginx" "-t"
  ;;
  "deploy_config")
    # 上送配置文件
    remoteSftp $config/config.yaml /tmp/config.yaml
    # 转移至项目目录
    remoteShell "mkdir -p $GO_PROJ_DIST && mv /tmp/config.yaml $GO_PROJ_DIST/config.yaml"
  ;;
  "deploy_vue")
    # vue构建与上送
    cd $basedir/../static && npm run build && cd -
    remoteShell "mkdir -p $STATIC_PROJ_DIST && rm -rf $STATIC_PROJ_DIST/*"
    remoteSftp "-r $basedir/../static/dist/*" $STATIC_PROJ_DIST
  ;;
  "deploy_linux_go")
    # 构建 go
    build linux
    # 不管目录是否存在仍创建
    remoteShell "mkdir -p $GO_PROJ_DIST/{bin,logs}"
    # 上送到指定目录
    PROJ_NAME=$PROJECT_NAME.`date +%s`
    remoteSftp $PROJECT $GO_PROJ_DIST/bin/$PROJ_NAME
    # 创建执行程序软连接
    remoteShell "ln -sf $GO_PROJ_DIST/bin/$PROJ_NAME $GO_PROJ_DIST/bin/$PROJECT_NAME"
  ;;
  "build_windows")
    build windows
  ;;
  "build_darwin")
    build darwin
  ;;
  *)
    echo "Usage: $0 $options"
  ;;
  esac
}

run $1