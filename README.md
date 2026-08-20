# otto-air

otto-air 是一个空气标准奥托循环（Otto cycle）核算器。用户输入压缩比 r、热容比 γ、进气状态（温度 T1、压力 p1）和单位质量加热量 q_in，后端计算四个过程端点的 p、v、T 状态，以及热效率 η、单位质量净功 w_net 和平均有效压力 MEP；并提供沿过程采样的 p–v 点列供绘制闭合曲线。循环按理想气体、定 γ 处理，气体常数 R = 287.0 J/(kg·K)，cv = R/(γ−1)。过程定义：1–2 等熵压缩、2–3 加热、3–4 等熵膨胀、4–1 等容放热，压缩比 r = v1/v2。加热段支持等容加热（奥托，默认）与等压加热（柴油对照，通过 mode 字段切换），两种模式共用同一套四段状态机，只切换加热过程类型。输入边界：r 必须大于 1、γ 必须大于 1、q_in 必须为正、进气温度与压力必须为正，违反任一条件均返回带 error 字段的 JSON 错误体（HTTP 400）。输出同时给出 η 的两种等价口径（1 − |q_out|/q_in 与 w_net/q_in），p–v 闭合曲线按数值积分给出的面积应与 w_net 一致。

## 启动

```bash
go run . -http :8080
```

打开 http://localhost:8080 进入交互页面，或直接用 curl 调用 API。

## API

两个求解接口都是 POST，请求与响应均为 JSON。

### POST /api/cycle

请求：

```json
{
  "r": 8,
  "gamma": 1.4,
  "intake": { "t": 300, "p": 101325 },
  "qin": 1200000,
  "mode": "otto"
}
```

- `r`：压缩比，必须大于 1。
- `gamma`：热容比 γ，必须大于 1。
- `intake`：进气状态，`t`（K）与 `p`（Pa）都必须为正。
- `qin`：单位质量加热量（J/kg），必须为正。
- `mode`：`"otto"`（等容加热，默认）或 `"diesel"`（等压加热对照），缺省为 otto。

响应含 `states`（点 1–4 的 p/v/T，SI 单位）、`eta`、`w_net`、`q_out`、`mep`、`cv`、`R` 与 `mode`。

### POST /api/pv

请求体与 `/api/cycle` 相同。响应含 `curve`（沿四个过程采样的 p–v 点列，闭合曲线）、`area`（数值积分面积，应逼近 `w_net`）与 `w_net`。页面用该点列在 SVG 上画 p–v 图。

### GET /api/examples

返回内置算例 `r8`（r=8、γ=1.4、T1=300 K、p1=101325 Pa、q_in=1.2 MJ/kg），供页面一键加载。

## 算例

- `example/r8.json`：r=8、γ=1.4，效率应等于 1−r^{1−γ} ≈ 0.5647，与 q_in 无关。

curl 验证：

```bash
curl -s -X POST http://localhost:8080/api/cycle \
  -H 'Content-Type: application/json' \
  --data-binary @example/r8.json
```

非法输入示例（r ≤ 1）：

```bash
curl -s -X POST http://localhost:8080/api/cycle \
  -H 'Content-Type: application/json' \
  -d '{"r":1,"gamma":1.4,"intake":{"t":300,"p":101325},"qin":1200000}'
```

返回 HTTP 400 与 `{"error":"compression ratio r must be greater than 1, got 1"}`。

## 测试

```bash
go test ./...
```

## 实现

领域内核在 `internal/`：`gas`（理想气体物性与校验）、`process`（过程类型、p–v 采样与数值积分）、`cycle`（循环求解）。`main.go` 只负责接线与静态页面托管。等熵关系、熵函数 s = cv ln T − R ln ρ 与端点计算共用 `gas` 包同一套物性，等熵段熵差应为零。
