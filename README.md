<!--
 * @Author: nijineko
 * @Date: 2025-06-10 13:54:48
 * @LastEditTime: 2025-06-10 23:18:12
 * @LastEditors: nijineko
 * @Description: 
 * @FilePath: \noa-encoder-json\README.md
-->
# Noa Encoder Json
JSON encoder for Noa Log. Supports encoding Log entries into JSON format for writing into log files.

# Installation
```bash
go get -u github.com/noa-log/noa-encoder-json
```

# Quick Start
```go
package main

import (
    "errors"
    "github.com/noa-log/noa"
    noaencoderjson "github.com/noa-log/noa-encoder-json"
)

func main() {
    // Create a new logger instance
    logger := noa.NewLog()
    // Set the encoder to the JSON encoder
    logger.SetEncoder(noaencoderjson.NewJSONEncoder(logger))

    // You can also set different encoders for print and write
    // logger.Encoder.Print = noa.NewTextEncoder(logger)
    // logger.Encoder.Write = noaencoderjson.NewJSONEncoder(logger)

    // Print Log
    logger.Info("Test", "This is an info message")
    logger.Error("Test", errors.New("an example error"))
}
```

## License
This project is open-sourced under the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0). Please comply with the terms when using it.