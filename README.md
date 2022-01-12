# Lightning 

## 前言
互联网越来越封闭，每当工作/学习中遇到什么问题，试图搜索排在最前面的都是 CSDN，往往点进去都是链中链，高质量的技术文章往往只有个人博客才会产出，但是由于大部分博主个人技术或者时间不充裕，往往被一些 SEO 做的非常好的网站挤到看不见的角落，藏在无人问津的黑暗中，类似 [travellings](https://github.com/volfclub/travellings) 这些项目的愿景其实非常好，网站和网站采用 One By One 的方式，可以让个人网站更大程度被挖掘出来，可是项目作者一直都是手动维护列表，新加入的网站平均要 1 个月才能被收录，并且在使用过程中，有的网站乱码或者被入侵挂上不良内容，已经会被用户访问到，于是我和 [lingyf](https://github.com/sxueck/lightning/stargazers) 同学发起了这个项目，希望能在这个黑暗中贡献出自己的 lightning

## 使用说明
项目分为几个大方向的开发阶段：
1. 复刻 travellings 项目 (已实现)
2. 加入个人网站信息提交页面，不需要走 issue (开发中)
3. 爬取计划成员的博文，将摘要分类汇总 (类似RSS)，用户对那篇文章感兴趣可以直接抵达目标博客进行与作者的交流
4. 火种计划，经过成员授权，抓取的博文将会以纯文本快照的形势存在服务器冷备，即使博客因为某些原因无法访问或者日后想要找回这些文章，可以直接调用快照
5. 共存，使用浏览器插件，当我们使用 Google / baidu 等搜索引擎，插件会自动插入与之关键词关联成员的相关博文

## 如何加入
当前项目还处于 Alpha 阶段，您暂时可以使用 https://dev.lighten.today  
建议插入在侧边栏或者顶部栏，例如:
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

## 其他文档
* [部署文档](documents/installation.md)
* [开发进度](documents/development.md)

## 致谢项目
* [友链接力](https://github.com/volfclub/travellings) : 使用并修改了项目的 index.html，根据 MPL 2.0 协议标准，需要指出对源码的修改，并沿用协议
