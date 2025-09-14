# gotris-term

[English](./README.md) | 简体中文

gotris-term 是一个演示项目，展示了如何为 gotris 实现自定义渲染器。gotris-term 使用基于文本的渲染，在终端中构建了一个基于 gotris 的迷你游戏，说明了如何用你自己的可视化输出扩展 gotris。

# 如何构建与游玩

## 构建可执行文件

在项目根目录下运行以下命令以构建 gotris-term 可执行文件：

```sh
go build -o gotris-term
```

这将在当前目录下生成名为 `gotris-term` 的可执行文件。

## 开始游戏

运行以下命令启动游戏：

```sh
./gotris-term
```

操作说明：

- 方向键：移动和旋转方块
- 空格键：硬降（Hard drop）
- Esc 或 Ctrl+C：退出游戏
