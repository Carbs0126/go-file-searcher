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
	// Terminal 除去help还能有多少行可以显示搜索结果
	SwitchScreenLines int
	// 当前显示的结果分组中，所有屏幕最高的高度
	MaxLineLength int
	// 当前状态显示几级菜单，0为不显示菜单
	CurrentMenuLevel int
	// Menu 显示时的高度
	MenuHeight int
	// Menu 显示时，左侧光标位置
	MenuLeftColumnIndex int
	// Menu 显示时，顶部光标位置
	MenuTopRowIndex int
	// 没有弹出menu时的cursor的位置
	MenuLevelOCursorIndex int
	// 弹出一级menu时的cursor的位置
	MenuLevel1CursorIndex int
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
		TerminalColumnNumber:  0,
		TerminalRowNumber:     0,
		SelectedLineIndex:     0,
		SelectedGroupIndex:    0,
		SwitchScreenLines:     0,
		MaxLineLength:         0,
		CurrentMenuLevel:      0,
		MenuHeight:            0,
		MenuLeftColumnIndex:   0,
		MenuTopRowIndex:       0,
		MenuLevelOCursorIndex: 0,
		MenuLevel1CursorIndex: 0,
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
