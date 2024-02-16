package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"regexp"
)

func main() {
	initStateData()
	getTerminalColumnsAndRows()
	getCommandState()
	getScreenLineNumber()
	if onlyPrintHelpInstructions() {
		printHelpInstructions()
		return
	}
	if gCommandState.Help {
		printHelpInstructions()
	}
	var re *regexp.Regexp = nil
	if len(gCommandState.SearchPattern) > 0 {
		re = regexp.MustCompile("(?i)" + gCommandState.SearchPattern)
	}
	printCurrentDirFiles(re)

	if len(gSearchData.FileDataArr) == 0 {
		fmt.Println("Search results are Empty.")
		return
	}
	clearPreviousNthLines(getSelectedGroupDisplayFileNamesLength(gSearchData.DisplayFileNamesInGroup))
	printCurrentLineWithSelectedDisplayName(gTerminalState.SelectedGroupIndex, gTerminalState.SelectedLineIndex)

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
			if gTerminalState.SelectedLineIndex < 0 {
				continue
			} else if gTerminalState.SelectedLineIndex == 0 {
				if 0 < gTerminalState.SelectedGroupIndex {
					showPreviousPage(gTerminalState.SelectedLineIndex, gTerminalState.MaxLineLength-1)
				}
			} else {
				clearCurrentLine()
				printCurrentLineWithUnselectedDisplayName(gTerminalState.SelectedGroupIndex, gTerminalState.SelectedLineIndex)
				gTerminalState.SelectedLineIndex = gTerminalState.SelectedLineIndex - 1
				if gTerminalState.SelectedLineIndex < getSelectedGroupDisplayFileNamesLength(gSearchData.DisplayFileNamesInGroup) {
					clearPreviousNthLines(1)
					printCurrentLineWithSelectedDisplayName(gTerminalState.SelectedGroupIndex, gTerminalState.SelectedLineIndex)
				}
			}
		} else if event.Key == keyboard.KeyArrowDown {
			if gTerminalState.SelectedLineIndex < getSelectedGroupDisplayFileNamesLength(gSearchData.DisplayFileNamesInGroup)-1 {
				// 当前页面
				clearCurrentLine()
				printCurrentLineWithUnselectedDisplayName(gTerminalState.SelectedGroupIndex, gTerminalState.SelectedLineIndex)
				gTerminalState.SelectedLineIndex = gTerminalState.SelectedLineIndex + 1
				clearNextLine()
				fmt.Print("\r")
				printCurrentLineWithSelectedDisplayName(gTerminalState.SelectedGroupIndex, gTerminalState.SelectedLineIndex)
			} else {
				if gTerminalState.SelectedGroupIndex < getGroupLength(gSearchData.DisplayFileNamesInGroup)-1 {
					// 如果有下一页，则进入下一页
					showNextPage(gTerminalState.SelectedLineIndex, 0)
				} else {
					// 没有下一页，到底了
					continue
				}
			}
		} else if event.Key == keyboard.KeyArrowRight {
			// 切屏
			if gTerminalState.SelectedGroupIndex >= getGroupLength(gSearchData.DisplayFileNamesInGroup)-1 {
				continue
			}
			showNextPage(gTerminalState.SelectedLineIndex, gTerminalState.SelectedLineIndex)
		} else if event.Key == keyboard.KeyArrowLeft {
			// 切屏
			if gTerminalState.SelectedGroupIndex <= 0 {
				continue
			}
			showPreviousPage(gTerminalState.SelectedLineIndex, gTerminalState.SelectedLineIndex)
		} else if (event.Key == keyboard.KeyEsc) || (event.Key == keyboard.KeyCtrlC) {
			clearCurrentLine()
			printCurrentLineWithUnselectedDisplayName(gTerminalState.SelectedGroupIndex, gTerminalState.SelectedLineIndex)
			fmt.Print("\r")
			clearNextNthLine(gTerminalState.MaxLineLength - gTerminalState.SelectedLineIndex)
			break
		} else if event.Key == keyboard.KeyEnter {
			// 打开文件
			ret := selectCurrentFile()
			if ret == 0 {
				clearNextNthLine(gTerminalState.MaxLineLength - gTerminalState.SelectedLineIndex)
				break
			}
		} else if event.Key == keyboard.KeySpace {
			// 打开文件的父目录
			ret := selectCurrentFilesParentDir()
			if ret == 0 {
				clearNextNthLine(gTerminalState.MaxLineLength - gTerminalState.SelectedLineIndex)
				break
			}
		}
	}
}

func showNextPage(oldSelectedLineIndex int, newSelectedLineIndex int) {
	gTerminalState.SelectedGroupIndex = gTerminalState.SelectedGroupIndex + 1
	// 光标先移动到最上面
	moveCursorToPreviousNthLines(oldSelectedLineIndex)
	// 遍历并清空屏幕
	for i := 0; i < gTerminalState.MaxLineLength; i++ {
		clearCurrentLine()
		fmt.Print("\r")
		if i < gTerminalState.MaxLineLength-1 {
			moveCursorToNextNthLines(1)
		}
	}
	// 光标移动到最上面
	moveCursorToPreviousNthLines(gTerminalState.MaxLineLength - 1)
	// 遍历并打印当前group的屏幕内容
	currentGroupLinesLength := getSelectedGroupDisplayFileNamesLength(gSearchData.DisplayFileNamesInGroup)
	for i := 0; i < currentGroupLinesLength; i++ {
		clearCurrentLine()
		printCurrentLineWithUnselectedDisplayName(gTerminalState.SelectedGroupIndex, i)
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
	gTerminalState.SelectedLineIndex = newSelectedLineIndex
	clearCurrentLine()
	printCurrentLineWithSelectedDisplayName(gTerminalState.SelectedGroupIndex, gTerminalState.SelectedLineIndex)
}

func showPreviousPage(oldPageSelectedLineIndex int, newSelectedLineIndex int) {
	gTerminalState.SelectedGroupIndex = gTerminalState.SelectedGroupIndex - 1
	// 光标先移动到最上面
	moveCursorToPreviousNthLines(oldPageSelectedLineIndex)
	// 遍历并清空屏幕
	for i := 0; i < gTerminalState.MaxLineLength; i++ {
		clearCurrentLine()
		fmt.Print("\r")
		if i < gTerminalState.MaxLineLength-1 {
			moveCursorToNextNthLines(1)
		}
	}
	// 光标移动到最上面
	moveCursorToPreviousNthLines(gTerminalState.MaxLineLength - 1)
	// 遍历并打印当前group的屏幕内容
	currentGroupLinesLength := getSelectedGroupDisplayFileNamesLength(gSearchData.DisplayFileNamesInGroup)
	for i := 0; i < currentGroupLinesLength; i++ {
		clearCurrentLine()
		printCurrentLineWithUnselectedDisplayName(gTerminalState.SelectedGroupIndex, i)
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
	gTerminalState.SelectedLineIndex = newSelectedLineIndex
	clearCurrentLine()
	printCurrentLineWithSelectedDisplayName(gTerminalState.SelectedGroupIndex, gTerminalState.SelectedLineIndex)
}
