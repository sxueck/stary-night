# Lightning

## 依赖创建
```sql
CREATE TABLE "sites" ("id" integer,"name" varchar NOT NULL,"url" varchar NOT NULL,"author" varchar NOT NULL, "lastmod" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (id));

INSERT INTO "sites" ("id", "name", "url", "author", "lastmod") VALUES
('1', 'epimetheus', 'https://www.sxueck.com/', 'sxueck', '2022-01-11 08:13:23');
```

## TODO
* 对于已经添加到了名单中的网站做篡改检查，防止变成不良内容
* 对于参加了项目的博客，爬虫会抓出文章摘要，主动进行内容聚合，如果用户对文章感兴趣，将可以直接跳转到源站点
* 后端每天对列表博客进行存活检测，如果网站异常，将被记录进低SLA名单，长时间处于该名单将会被移名
* 用户端每次进行访问也会对目标网站进行存活检测
