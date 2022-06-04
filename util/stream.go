package util

import (
	"bufio"
	"compiler/lexer/define"
	"container/list"
	"fmt"
	"io"
)

// Stream 流处理对象 用于确定源代码的词法分析的各个词法单元的位置 扫描并且输出结果
type Stream struct {
	scanner    *bufio.Scanner
	queueCache *list.List
	endToken   string
	isEnd      bool
	line       int
	column     int
}

// NewStream 创建一个 流处理对象 
func NewStream(r io.Reader, et string) *Stream {
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanRunes)
	return &Stream{scanner: s, queueCache: list.New(), endToken: et, isEnd: false, line: 1, column: 1}
}

// GetLine 获取所在行 用于定位分析报错
func (s *Stream) GetLine() int {
	return s.line
}

// GetColumn 获取所在列 用于定位分析报错
func (s *Stream) GetColumn() int {
	return s.column
}

// Next 获取下一个词法单元 或者流对象供给分析
func (s *Stream) Next() string {
	char := ""
	if s.queueCache.Len() != 0 {
		e := s.queueCache.Front()
		char = s.queueCache.Remove(e).(string)
	} else if s.scanner.Scan() {
		char = s.scanner.Text()
	} else {
		s.isEnd = true

		char = s.endToken
	}
	if define.IsNewLine(char) {
		s.line += 1
		s.column = 0
	}
	s.column += 1
	return char
}

// HasNext 是否还有下一个词法单元
func (s *Stream) HasNext() bool {
	if s.queueCache.Len() != 0 {
		return true
	}

	if s.scanner.Scan() {
		s.queueCache.PushBack(s.scanner.Text())
		return true
	}

	if !s.isEnd {
		return true
	}

	return false
}

// Peek 队列最前面的词法单元 取出并且返回
func (s *Stream) Peek() string {
	if s.queueCache.Len() != 0 {
		return s.queueCache.Front().Value.(string)
	}

	if s.scanner.Scan() {
		e := s.scanner.Text()
		s.queueCache.PushBack(e)
		return e
	}

	return s.endToken
}

// PutBack 将一个词法单元放回流对象
func (s *Stream) PutBack(e string) {
	s.queueCache.PushFront(e)
	s.column -= 1
}

// ClearFronts 跳过前面全部的流对象
func (s *Stream) ClearFronts(count int) {
	for i := 0; i < count; i++ {
		s.Next()
	}
}

// Print 输出流对象的内容 可用于调试和 debug
func (s *Stream) Print() {
	e := s.queueCache.Front()
	for {
		if e != nil {
			fmt.Print(e.Value)
			e = e.Next()
		} else {
			break
		}
	}
	fmt.Println("$")
}
