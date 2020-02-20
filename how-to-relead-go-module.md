# 如何优雅地发布 go module 模块
[TOC]

## 前言
截止到go1.13, go 官方推出的包管理工具go module已经发布三个版本了，网上也有很多文章介绍如何使用go module(推荐观看附录中Go夜读的视频和官方Wiki)，但是大部分都是讲如何引用别人的go module模块，鲜有提到如何发布自己的go module包的文章。本文将主要介绍如何“优雅”地发布自己的go module模块。

本文的代码和命令均在 ubuntu 19.10 中运行，go 的版本为 1.13.5


## pkg.go.dev简介
![](https://user-gold-cdn.xitu.io/2020/2/20/17061bd13d8f0c5e?w=1287&h=887&f=png&s=116940)
[go.dev](https://go.dev) 是go官方团队于2019年11月上线的集合go开发资源的网站，包括一些学习课程和一些go的案例，当然最重要的就提供了go的第三方包的检索功能。没错，他就是用来取代原来的[godoc.org](http://godoc.org)的，现在godoc.org上也有提示提醒用户迁移到pkg.go.dev。在这篇文章中，我们将把go module模块发布到pkg.go.dev。


BTW,在我写这篇文章的时候（2020.02.19），go官方刚好也宣布了go.dev不久将开源。

## 发布第一个版本
这次要发布的代码放在github，所以新建一个项目叫how-to-release-go-module
新增hello.go 文件
为hello.go添加两个方法和相关注释
```go
package pkg

import "fmt"

// Hello says hello.
func Hello() {
	fmt.Println("Hello go mod!")
}

// Bye says bye.
func Bye() {
	fmt.Println("Bye go mod!")
}
```

执行 `go mod init` 生成`go.mod`文件

```
go mod init
go: creating new go.mod: module github.com/YouEclipse/how-to-release-go-module

```
内容如下：
```
module github.com/YouEclipse/how-to-release-go-module

go 1.13
```
我们把代码push 到远端分支，看起来好像第一个版本就发布完毕了。我们打开`pkg.go.dev`搜索一下`github.com/YouEclipse/how-to-release-go-module`这个包，却返回未找到这个包。这是为何？

![](https://user-gold-cdn.xitu.io/2020/2/19/1705e054dfaca8bc?w=1209&h=196&f=png&s=28267)

其实在go.dev的[about](https://go.dev/about)中说的很清楚了，只有通过`proxy.golang.org`下载包的时候，module才会自动同步到pkg.go.dev。

> To add a package or module, simply fetch it from proxy.golang.org. 

但是，实际上，proxy.golang.org 国内基本上是无法访问的，如果我们使用 goproxy.cn,也一样能够同步，我没有研究goproxy的源码，但是我和goproxy的作者确认过，goproxy的推算行为会用到 proxy.golang.org，所以使用goproxy.cn作为代理也是可行的。

我们可以随便建一个项目，执行`go get -u github.com/YouEclipse/how-to-release-go-module`，因为go1.13 go module 已经是默认打开的，所以会默认通过proxy.golang.org拉取。如果不确定是否配置go proxy，可以执行`go env`和`go env -w` 命令查看和修改。

```
go get -u github.com/YouEclipse/how-to-release-go-module
go: finding github.com/YouEclipse/how-to-release-go-module latest
go: downloading github.com/YouEclipse/how-to-release-go-module v0.0.0-20200219150525-4f41ffd1dd8f
go: extracting github.com/YouEclipse/how-to-release-go-module v0.0.0-20200219150525-4f41ffd1dd8f
```
当我们成功拉取后，可以在pkg.go.dev再次搜索（具体可能需要等一段时间,大约是十分钟到半小时的样子），这时候我们可以看到搜索结果了

![](https://user-gold-cdn.xitu.io/2020/2/20/1705e6df2e2bc420?w=1239&h=915&f=png&s=102615)
![](https://user-gold-cdn.xitu.io/2020/2/20/1705e772303ec625?w=1172&h=306&f=png&s=50042)
看起来似乎我们的第一次发布大功告成了。我们看看pkg.go.dev包含了的什么信息：

- 版本号：由go module自动生成
- 发布时间
- 开源协议
- commit Hash
- tag:latest
- Overview(概览)
  - 包名
  - 源码地址
  - README
- Doc: godoc文档
- Subdirectories: 子目录
- Versions:已经发布过的版本
- Imports:引用的包
- Imported By: 引用此包的moudule
- Licenses：开源许可证
  
  ...

到了这里我们会发现，godoc 和 module的README 都没有正常显示,提示
`“Doc” not displayed due to license restrictions.`和`README not displayed due to license restrictions`，是说我们的包没有开源许可证，所以无法显示。对于这一点，网上有资料提到go官方团队和他们的律师讨论过才做的这个决定，这也是可以理解的。


pkg.go.dev 支持的证书可以在[https://pkg.go.dev/license-policy](https://pkg.go.dev/license-policy)查看，我们只要选择合适的开源协议证书添加到项目中即可（很遗憾，WTFPL并不在支持的开源协议中）。这里我们选择Apache2.0协议，添加到项目中，并push 到远端分支。

在等待一段时间（这里我等了大约30分钟）pkg.go.dev 更新后，我们会发现README和doc都可以正常显示了。这里生成的doc和godoc.org没太大的区别，都是基于代码和注释生成的，网上也有很多关于生成godoc最佳实践的文章，这里不做赘述。
![](https://user-gold-cdn.xitu.io/2020/2/20/1705e9f5be34d8c8?w=1177&h=684&f=png&s=64660)
![](https://user-gold-cdn.xitu.io/2020/2/20/1705e9ff0b16704c?w=1189&h=722&f=png&s=95971)

至此，发布的module包有godoc文档，有开源许可证，看起来是这么个样子，第一个版本至此就算发布完了。


## 发布新版本
其实我们在发布第一个版本的时候，为了更新license发布了两次，但是两次的版本都是v0.0.0,这么看起来似乎和go modules版本化的理念背道而驰。go module 实际上是可以通过tag来发布版本的。当我们需要发布新版本时，对应的，我们需要使用`git tag`为这个版本打上标签。假设我们发布的下个版本是v1.0.0：
```
git tag v1.0.0
```
打上标签后push到远端分支，等待一段时间，我们就可以在pkg.go.dev上看到我们发布的v1.0.0版本了。

![](https://user-gold-cdn.xitu.io/2020/2/20/170625bedef75e20?w=1169&h=322&f=png&s=46261)

（这里有些奇怪，之前的v0.0.0消失了，不知道什么原因）

我们执行`go get -u github.com/YouEclipse/how-to-release-go-module` 即可获取最新发布的版本
```shezhi
go get -u github.com/YouEclipse/how-to-release-go-module
go: finding github.com/YouEclipse/how-to-release-go-module v1.0.0
go: downloading github.com/YouEclipse/how-to-release-go-module v1.0.0
go: extracting github.com/YouEclipse/how-to-release-go-module v1.0.0
```



## 发布breaking changes
在早期没有go module时，假设我们引用的第三方包的做了breaking changes，API 发生改变，在跑CI或者重新拉取第三方包后，代码将会编译失败。我印象比较深刻的是在2018年左右，go.uuid的API发生了breaking changes，将原来没有返回err的函数返回了err，而当时我们没有任何包管理，都是在docker 镜像更新时通过go get 拉取，这就导致当时我们的CI都跑失败了。

go官方其实有关于版本控制的最佳实践，叫做[Semantic Import Version](https://github.com/golang/go/wiki/Modules)，即语义化版本。 关于语义化版本的说明，附录中官方的Wiki也有介绍，这里我们按照官方的最佳实践执行。
这里强调一下Go 官方对于语义化版本的一个基本原则：
> If an old package and a new package have the same import path, the new package must be backwards compatible with the old package."
> 
> 如果旧软件包和新软件包具有相同的导入路径，则新软件包必须与旧软件包向后兼容


所以Go官方给了两个方案，针对进行大版本升级和breaking changes：
 - Major branch 
    即通过创建version分支和tag进行版本升级
 - Major subdirectory
    即通过创建version子目录来区分不同版本

这里我们只使用Major branch方案来发布，因为第二种方案看起来很奇怪，而且似乎背离了VCS的意义，所有的版本代码居然都在一块，个人不是很推荐。

我们修改原来的代码，并做出breaking changes.

```go
package pkg

import "fmt"

// Hello says hello.
func Hello() error {
	fmt.Println("Hello go mod!")
	return nil
}

// Bye says bye.
func Bye() error {
	fmt.Println("Bye go mod!")
	return nil
}

```
并将import path 改为 v2
``` 
module github.com/YouEclipse/how-to-release-go-module/v2

go 1.13

```
然后提交代码，并创建tag
```
git tag v2.0.0
git push --tags
git push master
```
在又等待了一段时间后(pkg.go.dev更新确实是有点慢)，可以看到v2版本已经发布了，

![](https://user-gold-cdn.xitu.io/2020/2/20/17062c05965928a9?w=1160&h=465&f=png&s=64765)

这时候我们尝试在之前的测试项目拉取更新
```
go get -u github.com/YouEclipse/how-to-release-go-module

```
可以看到go.mod中显示我们使用的依然是v1.0.0,显然，v1.0.0版本的用户并没有受到breaking changes的影响。

如果我们要使用v2.0.0版本，修改go get的路径为`github.com/YouEclipse/how-to-release-go-module/v2`即可。

至此，我们的breaking changes版本的发布也完成了。


## 添加 go dev badge
大部分的开源的项目我们都可以在README中看到各种小图标，标识着项目的各种状态，一般称之为badge。在pkg.go.dev之前，大部分的go项目都会添加godoc.org的badge引导开发者们去godoc.org查看文档，但是既然使用了pkg.go.dev,我们自然就应该添加go.dev 的badge。更多的badge和相关设置可以在[shields.io](https://shields.io/)查看。

添加badge的方法和markdow添加图片的方法一样，只要替换项目在pkg.go.dev的路径即可
```markdown
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/url/of/your-module)
```
比如`github.com/YouEclipse/how-to-release-go-module/v2` 这个项目就可以设置成
```markdown
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/YouEclipse/how-to-release-go-module?tab=doc) 
```
效果如下

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](github.com/YouEclipse/how-to-release-go-module)

这样，我们的go module 模块看起来就很完美了。

## 结语
本文从一个简单的例子基本覆盖了发布go module 模块到pkg.go.dev可能会遇到的场景，希望能给阅读此文章的开发者提供帮助。




## 附录
1 [Go modules Wiki](https://github.com/golang/go/wiki/Modules)

2 [Go夜读第 61 期 Go Modules、Go Module Proxy 和 goproxy.cn](https://reading.developerlearning.cn/reading/61-2019-09-26-go-module-goproxy-cn/)


3 [go.dev: serve status badge similar to godoc.org](https://github.com/golang/go/issues/36982)

4 [示例项目 how-to-release-go-module](https://github.com/YouEclipse/how-to-release-go-module)

5 [Go module机制下升级major版本号的实践](https://tonybai.com/2019/06/03/the-practice-of-upgrading-major-version-under-go-module/)