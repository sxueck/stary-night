# Lightning

互联网越来越封闭，每当工作 / 学习中遇到什么问题，试图搜索排在最前面的 CSDN，点进去都是链中链，高质量的技术文章有的时候真的只有个人博客才会产出，但是由于大部分博主个人技术或者时间不充裕，往往被一些 SEO 做的非常好的网站挤到看不见的角落，藏在无人问津的黑暗中，类似 [travellings](https://github.com/volfclub/travellings) 这些项目的愿景其实非常好，网站和网站采用 One By One 的方式，可以让个人网站更大程度被挖掘出来，可是项目作者一直都是手动维护列表，新加入的网站平均要一个月才能被收录，并且在使用过程中，有的网站乱码或者被入侵挂上不良内容，已经会被用户访问到，于是我和 [lingyf](https://github.com/lingyf) 同学发起了这个项目，就是希望能在这个黑暗中贡献出属于自己的 lightning。

## 如何加入

当前项目还处于 Alpha 阶段，您暂时可以使用 https://dev.lighten.today；建议插入在侧边栏或者顶部栏，例如:
```yaml
menu:
  main:
    - identifier: AboutMe
      name: About
      url: /aboutme/
      weight: 10
    - identifier: ags
      name: Tags
      url: /tags/
      weight: 20
    - identifier: Friends
      name: <ls -al friends/*>
      weight: 30
      url: /friends/
    - identifier: lookstar
      name: 探星
      weight: 5
      url: https://dev.lighten.today
```
您可以使用 `探星` 或 `点灯` 等一切您认为合适的词汇

## 开发 & 功能

项目分为几个大方向的开发阶段：

1. 复刻 travellings 项目 (已实现)
2. 加入个人网站信息提交页面，不需要走 issue (开发中)
3. 爬取计划成员的博文，将摘要分类汇总 (类似 RSS)，用户对那篇文章感兴趣可以直接抵达目标博客进行与作者的交流
4. 火种计划，经过成员授权，抓取的博文将会以纯文本快照的形势存在服务器冷备，即使博客因为某些原因无法访问或者日后想要找回这些文章，可以直接调用快照
5. 共存，使用浏览器插件，当我们使用 Google / Baidu 等搜索引擎，插件会自动插入与之关键词关联成员的相关博文

## 进度 & 文档

* [部署文档](docs/installation.md)
* [开发进度](docs/development.md)

## 协议 & 说明

我们充分保障您的自由权利，但相应的，我们也需要遵循一些君子协议：

* 您的网站可以是自建站或者托管站 (简书 / 博客园等)
* 我们推荐您使用 Sitemap 减低来自 Lightning 的爬虫难度
* 您可以畅所欲言，当然意义上的 "不良内容" 除外
* 您随时能退出计划，也能随时销毁来自您已经被收录的文章，这是您理所应当的权利
* 针对项目的发展或者需求，您可以在 Issue 提出想法，甚至可以参与开发进程中

## 致谢 & 许可

* [友链接力](https://github.com/volfclub/travellings) : 使用并修改了项目的 index.html，根据 MPL 2.0 协议标准，需要指出对源码的修改并沿用协议，详情参考 [LICENSE](LICENSE)