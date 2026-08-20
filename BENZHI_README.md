# otto-air

otto-air 是一个空气标准奥托循环（Otto cycle）核算器。用户输入压缩比 r、热容比 γ、进气状态（温度 T1、压力 p1）和单位质量加热量 q_in，后端计算四个过程端点的 p、v、T 状态，以及热效率 η、单位质量净功 w_net 和平均有效压力 MEP；并提供沿过程采样的 p–v 点列供绘制闭合曲线。

## 构建 / 运行 / 测试

```text
go build ./...     # 编译
go run . -http :8080
go test ./...      # 测试
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
