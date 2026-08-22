# chvorinov-t：铸件凝固时间 Chvorinov 核算命令行工具

按 Chvorinov 规则由铸件体积 V、表面积 A、模具常数 C 与指数 n 计算凝固时间 tf = C·(V/A)ⁿ，并比较板/柱/球/立方体各形状的模数。

## 构建 / 运行 / 测试

```text
go build ./...
go run . freeze example/steel-cube.json
go test ./...
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
