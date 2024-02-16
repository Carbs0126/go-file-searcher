package main

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
	// 所有文件名称，第一维是屏幕
	DisplayFileNamesInGroup [][]string
	// 所有文件名称
	DisplayFileNames []string
	// 文件真实名称，从当前目录开始
	RealFileNames []string
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
		DisplayFileNamesInGroup: make([][]string, 0, 2),
		DisplayFileNames:        make([]string, 0, 16),
		RealFileNames:           make([]string, 0, 16),
	}
}
