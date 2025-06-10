/*
 * @Author: nijineko
 * @Date: 2025-06-10 12:57:28
 * @LastEditTime: 2025-06-10 23:02:39
 * @LastEditors: nijineko
 * @Description: noa json log encoder
 * @FilePath: \noa-encoder-json\jsonEncoder.go
 */
package noaencoderjson

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/noa-log/colorize"
	"github.com/noa-log/noa"
	"github.com/noa-log/noa/encoder"
	"github.com/noa-log/noa/errors"
)

// json log format
type JSONLog struct {
	Time       int64  `json:"time"`
	TimeFormat string `json:"time_format"`
	Level      int    `json:"level"`
	Source     string `json:"source"`
	Data       []any  `json:"data"`
}

// error log stack frame format
type ErrorLogStackFrame struct {
	Function     string `json:"function"`
	File         string `json:"file"`
	Line         int    `json:"line"`
	PackageName  string `json:"package_name"`
	FunctionName string `json:"function_name"`
}

// error log format
type ErrorLog struct {
	Message     string               `json:"message"`
	StackFrames []ErrorLogStackFrame `json:"stack_frames,omitempty"`
}

// JSONEncoder struct
type JSONEncoder struct {
	Log *noa.LogConfig
}

/**
 * @description: Create a new JSONEncoder instance
 * @param {noa.LogConfig} Log noa log instance
 * @return {*JSONEncoder} JSONEncoder instance
 */
func NewJSONEncoder(Log *noa.LogConfig) *JSONEncoder {
	return &JSONEncoder{
		Log: Log,
	}
}

// print json log data
func (je *JSONEncoder) Print(c *encoder.Context) {
	JSONData, err := marshalJSON(je.Log, c)
	if err != nil {
		panic(err)
	}

	// cache json data in context
	c.Set("PrintData", JSONData)

	fmt.Printf("%s\n", JSONData)
}

// return file extension for the encoded file
func (je *JSONEncoder) WriteFileExtension() string {
	return ".json"
}

// write log data to file
func (je *JSONEncoder) Write(FileHandle *os.File, c *encoder.Context) error {
	JSONData := c.Get("PrintData")
	JSONDataBytes, ok := JSONData.([]byte)
	if JSONData == nil || !ok {
		JSON, err := marshalJSON(je.Log, c)
		if err != nil {
			return err
		}

		JSONDataBytes = JSON
	}

	// add newline character
	JSONDataBytes = append(JSONDataBytes, '\n')

	if _, err := FileHandle.Write(JSONDataBytes); err != nil {
		return err
	}

	return nil
}

/**
 * @description: Get file extension for the encoded file
 * @return {string} file extension
 */
func (e *JSONEncoder) FileExtension() string {
	return ".json"
}

/**
 * @description: Marshal JSON log entry
 * @param {*noa.LogConfig} l noa log instance
 * @param {*encoder.Context} c encoder context
 * @return {[]byte} JSON encoded log entry
 * @return {error} error
 */
func marshalJSON(l *noa.LogConfig, c *encoder.Context) ([]byte, error) {
	LogEntry := JSONLog{
		Time:       c.Time.Unix(),
		TimeFormat: c.Time.Format(l.TimeFormat),
		Level:      c.Level,
		Source:     c.Source,
	}

	for Index, Data := range c.Data {
		if l.RemoveColor {
			if DataStr, ok := Data.(string); ok {
				Data = colorize.Remove(DataStr)
				c.Data[Index] = Data
			}
		}

		if l.Errors.StackTrace {
			// if data is error wrap
			if DataError, ok := Data.(*errors.Error); ok {
				var ErrorLogStackFrames []ErrorLogStackFrame
				for _, Frames := range DataError.StackFrames() {
					ErrorLogStackFrames = append(ErrorLogStackFrames, ErrorLogStackFrame{
						Function:     Frames.Function,
						File:         Frames.File,
						Line:         Frames.Line,
						PackageName:  Frames.PackageName,
						FunctionName: Frames.FunctionName,
					})
				}

				LogEntry.Data = append(LogEntry.Data, ErrorLog{
					Message:     DataError.Error(),
					StackFrames: ErrorLogStackFrames,
				})

				continue
			}
		}
		// skip newline characters
		if _, ok := Data.(string); ok {
			if Data == "\n" || Data == "\r\n" {
				continue
			}
		}

		LogEntry.Data = append(LogEntry.Data, Data)
	}

	// json encode
	return json.Marshal(LogEntry)
}
