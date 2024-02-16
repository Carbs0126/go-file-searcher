package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"os"
	"regexp"
)

// 命令行屏幕列数
var gTerminalColumnNumber = 0

// 命令行屏幕行数
var gTerminalRowNumber = 0

// 当前屏幕光标所在行索引
var gSelectedLineIndex = 0

// 当前屏幕索引
var gSelectedGroupIndex = 0

// 所有文件名称，第一维是屏幕
var gFileDisplayNamesGroup = make([][]string, 0, 2)

// 所有文件名称
var gFileDisplayNames = make([]string, 0, 16)

// 文件真实名称，从当前目录开始
var gFileRealNames = make([]string, 0, 16)

// 一屏幕多少行显示文件
var gSwitchScreenLines = 0

// 所有屏幕最高的高度
var gMaxLineLength = 0

// 是否包含 -r 命令，递归遍历当前文件夹
var gContainRecursive = false

func main() {
	gTerminalColumnNumber, gTerminalRowNumber = getTerminalColumnsAndRows()
	args := os.Args[1:]
	argsLength := len(args)
	if argsLength == 0 {
		gSwitchScreenLines = gTerminalRowNumber - 2
		printCurrentDirFiles(nil)
	} else if argsLength == 1 && isArgHelp(args[0]) {
		printInstructions()
		return
	} else {
		gContainRecursive = containArgRecursive(args)
		gSwitchScreenLines = gTerminalRowNumber - 2
		if containArgHelp(args) {
			printInstructions()
			gSwitchScreenLines = gTerminalRowNumber - 2 - len(gInstructions)
		}
		pattern := getSearchPattern(args)
		re := regexp.MustCompile("(?i)" + pattern)
		printCurrentDirFiles(re)
	}
	if len(gFileDisplayNames) == 0 {
		fmt.Println("Search results are Empty.")
		return
	}
	clearPreviousNthLines(getSelectedGroupDisplayFileNamesLength())
	printCurrentLineWithSelectedDisplayName(gSelectedGroupIndex, gSelectedLineIndex)

	err := keyboard.Open()
	if err != nil {
		fmt.Println("Error occurred when opening keyboard:", err)
		return
	}
	defer keyboard.Close()

	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		fmt.Println("Error occurred when getting keys:", err)
		return
	}
	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}
		if event.Key == keyboard.KeyArrowUp {
			if gSelectedLineIndex < 0 {
				continue
			} else if gSelectedLineIndex == 0 {
				if 0 < gSelectedGroupIndex {
					showPreviousPage(gSelectedLineIndex, gMaxLineLength-1)
				}
			} else {
				clearCurrentLine()
				printCurrentLineWithUnselectedDisplayName(gSelectedGroupIndex, gSelectedLineIndex)
				gSelectedLineIndex = gSelectedLineIndex - 1
				if gSelectedLineIndex < getSelectedGroupDisplayFileNamesLength() {
					clearPreviousNthLines(1)
					printCurrentLineWithSelectedDisplayName(gSelectedGroupIndex, gSelectedLineIndex)
				}
			}
		} else if event.Key == keyboard.KeyArrowDown {
			if gSelectedLineIndex < getSelectedGroupDisplayFileNamesLength()-1 {
				// 当前页面
				clearCurrentLine()
				printCurrentLineWithUnselectedDisplayName(gSelectedGroupIndex, gSelectedLineIndex)
				gSelectedLineIndex = gSelectedLineIndex + 1
				clearNextLine()
				fmt.Print("\r")
				printCurrentLineWithSelectedDisplayName(gSelectedGroupIndex, gSelectedLineIndex)
			} else {
				if gSelectedGroupIndex < getGroupLength()-1 {
					// 如果有下一页，则进入下一页
					showNextPage(gSelectedLineIndex, 0)
				} else {
					// 没有下一页，到底了
					continue
				}
			}
		} else if event.Key == keyboard.KeyArrowRight {
			// 切屏
			if gSelectedGroupIndex >= getGroupLength()-1 {
				continue
			}
			showNextPage(gSelectedLineIndex, gSelectedLineIndex)
		} else if event.Key == keyboard.KeyArrowLeft {
			// 切屏
			if gSelectedGroupIndex <= 0 {
				continue
			}
			showPreviousPage(gSelectedLineIndex, gSelectedLineIndex)
		} else if (event.Key == keyboard.KeyEsc) || (event.Key == keyboard.KeyCtrlC) {
			clearCurrentLine()
			printCurrentLineWithUnselectedDisplayName(gSelectedGroupIndex, gSelectedLineIndex)
			fmt.Print("\r")
			clearNextNthLine(gMaxLineLength - gSelectedLineIndex)
			break
		} else if event.Key == keyboard.KeyEnter {
			// 打开文件
			ret := selectCurrentFile()
			if ret == 0 {
				clearNextNthLine(gMaxLineLength - gSelectedLineIndex)
				break
			}
		} else if event.Key == keyboard.KeySpace {
			// 打开文件的父目录
			ret := selectCurrentFilesParentDir()
			if ret == 0 {
				clearNextNthLine(gMaxLineLength - gSelectedLineIndex)
				break
			}
		}
	}
}

func showNextPage(oldSelectedLineIndex int, newSelectedLineIndex int) {
	gSelectedGroupIndex = gSelectedGroupIndex + 1
	// 光标先移动到最上面
	moveCursorToPreviousNthLines(oldSelectedLineIndex)
	// 遍历并清空屏幕
	for i := 0; i < gMaxLineLength; i++ {
		clearCurrentLine()
		fmt.Print("\r")
		if i < gMaxLineLength-1 {
			moveCursorToNextNthLines(1)
		}
	}
	// 光标移动到最上面
	moveCursorToPreviousNthLines(gMaxLineLength - 1)
	// 遍历并打印当前group的屏幕内容
	currentGroupLinesLength := getSelectedGroupDisplayFileNamesLength()
	for i := 0; i < currentGroupLinesLength; i++ {
		clearCurrentLine()
		printCurrentLineWithUnselectedDisplayName(gSelectedGroupIndex, i)
		fmt.Print("\r")
		if i < currentGroupLinesLength-1 {
			moveCursorToNextNthLines(1)
		}
	}
	// 确保 newSelectedLineIndex 不超过当前页面内容
	if newSelectedLineIndex >= currentGroupLinesLength {
		newSelectedLineIndex = currentGroupLinesLength - 1
	}
	// 从当前内容的最后一行，定位到想要光标定位的行
	moveCursorToPreviousNthLines(currentGroupLinesLength - 1 - newSelectedLineIndex)
	gSelectedLineIndex = newSelectedLineIndex
	clearCurrentLine()
	printCurrentLineWithSelectedDisplayName(gSelectedGroupIndex, gSelectedLineIndex)
}

func showPreviousPage(oldPageSelectedLineIndex int, newSelectedLineIndex int) {
	gSelectedGroupIndex = gSelectedGroupIndex - 1
	// 光标先移动到最上面
	moveCursorToPreviousNthLines(oldPageSelectedLineIndex)
	// 遍历并清空屏幕
	for i := 0; i < gMaxLineLength; i++ {
		clearCurrentLine()
		fmt.Print("\r")
		if i < gMaxLineLength-1 {
			moveCursorToNextNthLines(1)
		}
	}
	// 光标移动到最上面
	moveCursorToPreviousNthLines(gMaxLineLength - 1)
	// 遍历并打印当前group的屏幕内容
	currentGroupLinesLength := getSelectedGroupDisplayFileNamesLength()
	for i := 0; i < currentGroupLinesLength; i++ {
		clearCurrentLine()
		printCurrentLineWithUnselectedDisplayName(gSelectedGroupIndex, i)
		fmt.Print("\r")
		if i < currentGroupLinesLength-1 {
			moveCursorToNextNthLines(1)
		}
	}
	// 确保 newSelectedLineIndex 不超过当前页面内容
	if newSelectedLineIndex >= currentGroupLinesLength {
		newSelectedLineIndex = currentGroupLinesLength - 1
	}
	// 从当前内容的最后一行，定位到想要光标定位的行
	moveCursorToPreviousNthLines(currentGroupLinesLength - 1 - newSelectedLineIndex)
	gSelectedLineIndex = newSelectedLineIndex
	clearCurrentLine()
	printCurrentLineWithSelectedDisplayName(gSelectedGroupIndex, gSelectedLineIndex)
}
