/*
 * @Author: nijineko
 * @Date: 2025-06-10 12:57:28
 * @LastEditTime: 2025-06-10 13:25:24
 * @LastEditors: nijineko
 * @Description: noa json log encoder
 * @FilePath: \noa-encoder-json\jsonEncoder_test.go
 */
package noaencoderjson

import (
	"errors"
	"testing"

	"github.com/noa-log/noa"
)

func TestJSONEncoder(t *testing.T) {
	Log := noa.NewLog()
	Log.Writer.Encoder = NewJSONEncoder(Log)

	Log.Debug("Test", "This is an info message", "key1", "value1", "key2", 123)
	Log.Info("Test", "This is an info message", "key1", "value1", "key2", 123)
	Log.Warning("Test", "This is a warning message")
	Log.Error("Test", errors.New("an example error"))
	Log.Fatal("Test", errors.New("a fatal error"))
}
