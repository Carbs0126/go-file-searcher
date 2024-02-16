package main

import (
	"time"
)

type TerminalState struct {
	// 命令行屏幕列数
	TerminalColumnNumber int
	// 命令行屏幕行数
	TerminalRowNumber int
	// 当前屏幕光标所在行索引
	SelectedLineIndex int
	// 当前屏幕索引
	SelectedGroupIndex int
	// 一屏幕多少行显示文件
	SwitchScreenLines int
	// 所有屏幕最高的高度
	MaxLineLength int
}

type CommandState struct {
	// 是否包含 help
	Help bool
	// 是否包含 recursive
	Recursive bool
	// 是否包含 time
	Time bool
	// 是否包含 search word
	SearchPattern string
}

type SearchData struct {
	// 文件信息
	FileDataArr []FileData
	// 所有文件名称，第一维是屏幕
	DisplayFileNamesInGroup [][]string
}

type FileData struct {
	DisplayFileName string
	FilePath        string
	Time            time.Time
}

var gTerminalState TerminalState
var gCommandState CommandState
var gSearchData SearchData

func initStateData() {
	gTerminalState = TerminalState{
		TerminalColumnNumber: 0,
		TerminalRowNumber:    0,
		SelectedLineIndex:    0,
		SelectedGroupIndex:   0,
		SwitchScreenLines:    0,
		MaxLineLength:        0,
	}
	gCommandState = CommandState{
		Help:          false,
		Recursive:     false,
		Time:          false,
		SearchPattern: "",
	}
	gSearchData = SearchData{
		FileDataArr:             make([]FileData, 0, 32),
		DisplayFileNamesInGroup: make([][]string, 0, 2),
	}
}

// 针对 FileData 的排序

type FileDataSlice []FileData

func (fd FileDataSlice) Len() int {
	return len(fd)
}

func (fd FileDataSlice) Less(i, j int) bool {
	return (fd[i].Time.UnixMilli() - fd[j].Time.UnixMilli()) > 0
}

func (fd FileDataSlice) Swap(i, j int) {
	fd[i], fd[j] = fd[j], fd[i]
}
