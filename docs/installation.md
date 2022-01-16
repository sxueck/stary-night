# Lightning 自部署文档

> 建议使用 Docker 进行部署

## 环境变量

* `HTTP_PORT` : 后台 Server 的 HTTP 端口
* `LISTEN_ADDRESS` : 后台 Server 的 HTTP 监听地址，建议保持默认参数，使用 Nginx 等服务进行反向代理
* `DB_NAME` : 数据库名称，使用 Sqlite3 作为数据库

## 依赖创建

使用 SQL 语句创建数据库

```shell
$ sqlite3 storage.db
```

```sql
CREATE TABLE "sites"
(
    "id"          integer PRIMARY KEY autoincrement,
    "name"        varchar  NOT NULL,
    "url"         varchar  NOT NULL,
    "author"      varchar  NOT NULL,
    "lastmod"     datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "contact"     varchar  NOT NULL,
    "description" varchar  NOT NULL
);
```

可以插入一条测试数据

```sql
INSERT INTO "sites" ("id", "name", "url", "author", "lastmod","contact","description")
VALUES ('1', 'epimetheus', 'https://www.sxueck.com/', 'sxueck', '2022-01-11 08:13:23','sxuecks@gmail.com','a sites');
```

## 直接编译

程序使用 Go 语言进行编写，由于 Go Modules 的特性，编译只需要注意代理，非常简单

```shell
$ export GOPROXY=https://goproxy.cn
$ go build .
$ ./starry-night
```

编译时请保证机器内存大于 512M (包括 Swap)，不然会导致 GCC 异常

## 使用 Docker 进行编译
```shell
$ docker build -t starry-night .
```

## Help

### 我内存太小了怎么办

直接使用虚拟内存即可，这里的 `SWAP_SIZE` 换成真实内存的 2 倍，例如真实内存为 1G，这里应该填 2G

```shell
# fallocate -l 2g swap
$ fallocate -l SWAP_SIZE swap
$ mkswap swap
$ swapon swap
```