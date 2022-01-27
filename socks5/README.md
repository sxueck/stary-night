# Socks5 Proxy

**a simple custom proxy service**

特点:
* 自定义日志格式
* 记录详细的网络情况
* 自定义白名单
* 部署可控简单

环境变量:
* SOCKS_PORT： 代理端口，默认 `13030`
* WHITE_ADDR：白名单模式，不设置的话将无法对外提供访问，格式 `xx.xx.xx.xx,zz.zz.zz.zz,xx.xx.xx.xx`
