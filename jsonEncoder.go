/*
 * @Author: nijineko
 * @Date: 2025-06-10 12:57:28
 * @LastEditTime: 2025-06-10 13:40:17
 * @LastEditors: nijineko
 * @Description: noa json log encoder
 * @FilePath: \noa-encoder-json\jsonEncoder.go
 */
package noaencoderjson

import (
	"encoding/json"
	"os"

	"github.com/noa-log/noa"
	"github.com/noa-log/noa/errors"
	"github.com/noa-log/noa/tools/output"
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

/**
 * @description: Get file extension for the encoded file
 * @return {string} file extension
 */
func (e *JSONEncoder) FileExtension() string {
	return ".json"
}

/**
 * @description: JSON encoder write method
 * @param {*os.File} FileHandle file handle
 * @param {[]any} PrintData data to encode
 * @return {error} error if any
 */
func (e *JSONEncoder) Write(FileHandle *os.File, PrintData []any) error {
	// unwrap the print data
	Time, Level, Source, PrintData := output.UnwrapPrintData(e.Log, PrintData)

	LogEntry := JSONLog{
		Time:       Time.Unix(),
		TimeFormat: Time.Format(e.Log.TimeFormat),
		Level:      Level,
		Source:     Source,
	}

	LastErrorIndex := -1
	for Index, Data := range PrintData {
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
			LastErrorIndex = Index

			continue
		}
		if _, ok := Data.(string); ok {
			// if previous line is error log, skip 2 lines
			if LastErrorIndex != -1 && (LastErrorIndex == Index-1 || LastErrorIndex == Index-2) {
				continue
			}

			// skip newline characters
			if Data == "\n" || Data == "\r\n" {
				continue
			}
		}

		LogEntry.Data = append(LogEntry.Data, Data)
	}

	// json encode
	JSONData, err := json.Marshal(LogEntry)
	if err != nil {
		return err
	}

	// write to file
	if _, err := FileHandle.Write(append(JSONData, '\n')); err != nil {
		return err
	}

	return nil
}
