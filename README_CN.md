# Noa Encoder Json
Noa Log 的 JSON 编码器。支持将 Log 项编码为 JSON 格式写入日志文件。

# 安装
```bash
go get -u github.com/noa-log/noa-encoder-json
```

# 快速开始
```go
package main

import (
    "errors"
    "github.com/noa-log/noa"
    noaencoderjson "github.com/noa-log/noa-encoder-json"
)

func main() {
    // 创建一个新的日志实例
    logger := noa.NewLog()
    // 设置编码器为 JSON 编码器
    logger.SetEncoder(noaencoderjson.NewJSONEncoder(logger))

    // 你也可以为打印和写入设置不同的编码器
    // logger.Encoder.Print = noa.NewTextEncoder(logger)
    // logger.Encoder.Write = noaencoderjson.NewJSONEncoder(logger)

    // 记录日志
    logger.Info("Test", "This is an info message")
    logger.Error("Test", errors.New("an example error"))
}
```

## 许可
本项目基于[Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0)协议开源。使用时请遵守协议的条款。