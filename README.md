# image-sync

docker镜像同步工具

该工具可以实现跨网络的镜像仓库之间同步镜像，可以快速复制镜像layer二进制

# 快速开始

## 安装
下载release的客户端二进制工具

## 使用

### 同步镜像
```bash
image-sync-tool sync <镜像> \
--src-registry-username=<镜像仓库账号> \
--src-registry-password=<镜像仓库密码> \
--server=<服务端地址> \
--username=<服务端账号> \
--password=<服务端密码> \
--dest-deploy-host=<部署主机地址> \
--dest-deploy-port=<部署主机ip> \
--dest-deploy-compose-file=<部署主机docker compose文件目录>
```

示例:

```bash
image-sync-tool sync harbor.hy-zw.com/park/web-park-admin:20240705022155 \
--src-registry-username=**** \
--src-registry-password="*****" \
--server=******* \
--username=***** \
--password=****** \
--dest-deploy-host=10.226.22.6 \
--dest-deploy-port=36000 \
--dest-deploy-compose-file=/data/docker-compose.yml
```