# genx

递归地查找给定目录下所有 `//go:generate` 指令，并对其所在目录执行 `go generate` 命令。

## 前言

经过 `go` 语言社区在没有宏和泛型的条件下的不断发展壮大，对于如何提高 `go` 开发效率，社区及各路开发者已经给出了自己的答案，即使用 `go generate` 代码生成。相对知名的库和框架都多多少少涉及代码生成，例如：

- [ent](https://entgo.io/)
- [grpc](https://grpc.io/)
- [go-zero](https://github.com/zeromicro/go-zero)
- [easyjson](https://github.com/mailru/easyjson)

基本上，`go` 语言开发的现状就是，项目规模稍微大一些，就开始离不开代码生成（除了各种库提供的代码生成外，还包括项目开发者自己提供的代码生成功能）；而这些代码生成指令 `//go:generate` 散落在项目的各个角落，如果需要一次性执行整个项目所有的 `//go:generate` 指令，只能一个目录一个目录去 `cd && go generate`，好一点的方式是写在 `.sh` 脚本或者 `makefile` 中，但仍然需要开发者手动维护，每次新增一条 `//go:generate` 指令，都需要更新 `.sh` 或 `makefile`。

为了简化在项目中执行 `go generate` 命令的方式，初步实现了第一版 `genx`，`genx` 会递归地查找给定目录下所有 `//go:generate` 指令，并在其所在目录执行 `go generate` 命令。通常而言，我们只需要在项目根目录执行一次 `genx` 命令即可完成整个项目的代码生成操作，而无需逐个目录去找需要进行代码生成的文件目录路径，也无需 `.sh` 和 `makefile`。同时，`genx` 还会打印出已执行的 `//go:generate` 指令，以便开发者能获知哪些代码生成指令已被执行，执行的命令是什么，所在的 `.go` 文件是哪个，在文件中的哪一行。

## 安装

```shell
go install github.com/x5iu/genx
```

## 使用

基本的使用方式为，在需要执行 `go generate` 的项目根目录下执行：

```shell
genx
genx .
genx service/wechat
```

`genx` 接收若干个位置参数，这些位置参数应为需要执行 `go generate` 的目录路径（仅允许目录路径作为参数，而不允许为文件路径，如果传入文件路径，`genx` 将会返回一个错误），如果不传入任何参数，`genx` 将会默认从当前工作路径开始执行 `go generate`。

`genx` 将从给定的目录路径开始，递归地扫描该目录及目录下所有子目录中的 `.go` 文件，识别其中内嵌的 `//go:generate` 指令，并在该目录下执行 `go generate` 命令；如果同一个目录下包含多个 `//go:generate` 指令，那么这个目录也只会被执行一次 `go generate`，而不会多次重复执行。

`genx` 还会列出已执行的 `//go:generate` 指令，如果你希望只查看 `//go:generate` 指令而不实际执行 `go generate` 命令，可以添加 `-l/--list` 参数：

```shell
genx -l .
genx --list .
```



