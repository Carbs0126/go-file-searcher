package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"regexp"
	"runtime"
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
		} else if (event.Key == keyboard.KeyEsc) || (event.Rune == 'Q') || (event.Rune == 'q') {
			if gTerminalState.CurrentMenuLevel > 0 {
				if gTerminalState.CurrentMenuLevel == 1 {
					dismissMenuLevel1()
				}
			} else {
				restoreCurrentLineWithSearchResult()
				clearNextNthLine(gTerminalState.MaxLineLength - gTerminalState.SelectedLineIndex)
				break
			}
		} else if event.Key == keyboard.KeyCtrlC {
			restoreCurrentLineWithSearchResult()
			clearNextNthLine(gTerminalState.MaxLineLength - gTerminalState.SelectedLineIndex)
			break
		} else if event.Key == keyboard.KeyEnter {
			if isMac() {
				// 打开文件
				ret := openCurrentFile()
				if ret == 0 {
					clearNextNthLine(gTerminalState.MaxLineLength - gTerminalState.SelectedLineIndex)
					break
				}
			} else {
				popMenu()
			}
		} else if (event.Rune == 'P') || (event.Rune == 'p') {
			if isMac() {
				// 打开文件的父目录
				ret := openCurrentFilesParentDir()
				if ret == 0 {
					clearNextNthLine(gTerminalState.MaxLineLength - gTerminalState.SelectedLineIndex)
					break
				}
			} else {
				popMenu()
			}
		} else if event.Key == keyboard.KeySpace {
			popMenu()
		}
	}
}

func restoreCurrentLineWithSearchResult() {
	clearCurrentLine()
	printCurrentLineWithUnselectedDisplayName(gTerminalState.SelectedGroupIndex, gTerminalState.SelectedLineIndex)
	fmt.Print("\r")
}

func isLinux() bool {
	return runtime.GOOS == "linux"
}

func isMac() bool {
	return runtime.GOOS == "darwin"
}

// 打开菜单对话框
func popMenu() {
	if gTerminalState.CurrentMenuLevel == 0 {
		gTerminalState.CurrentMenuLevel = 1
		popMenuLevel1()
	}
}

func popMenuLevel1() {
	// 菜单高度
	gTerminalState.MenuHeight = len(gMenu)
	// 菜单距离左侧 16 个字符
	gTerminalState.MenuLeftColumnIndex = 10
	// 菜单顶部
	gTerminalState.MenuTopRowIndex = 0
	if gTerminalState.MaxLineLength-gTerminalState.SelectedLineIndex > gTerminalState.MenuHeight {
		// 当前光标的下部分高度大于menu高度
		gTerminalState.MenuTopRowIndex = gTerminalState.SelectedLineIndex + 1
	} else if gTerminalState.SelectedLineIndex >= gTerminalState.MenuHeight {
		// 当前光标的下部分高度小于menu高度，且光标上部分高度大于menu高度
		gTerminalState.MenuTopRowIndex = gTerminalState.SelectedLineIndex - gTerminalState.MenuHeight
	} else if gTerminalState.MaxLineLength > gTerminalState.MenuHeight {
		// 居中
		gTerminalState.MenuTopRowIndex = (gTerminalState.MaxLineLength-1-gTerminalState.MenuHeight)/2 + 1
	} else {
		// todo 拓展高度？
	}
	if gTerminalState.SelectedLineIndex < gTerminalState.MenuTopRowIndex {
		// 下移光标
		moveCursorToNextNthLines(gTerminalState.MenuTopRowIndex - gTerminalState.SelectedLineIndex)
	} else {
		// 上移光标
		moveCursorToPreviousNthLines(gTerminalState.SelectedLineIndex - gTerminalState.MenuTopRowIndex)
	}
	gTerminalState.MenuLevelOCursorIndex = gTerminalState.SelectedLineIndex
	gTerminalState.SelectedLineIndex = gTerminalState.MenuTopRowIndex
	// 开始输出文字
	for i := 0; i < gTerminalState.MenuHeight; i++ {
		// 光标移动到合适的列位置
		jumpToColumnIndex(gTerminalState.MenuLeftColumnIndex)
		// 输出文字
		fmt.Print(gMenu[i])
		// 光标移动到下一行
		moveCursorToNextNthLines(1)
	}
	gTerminalState.SelectedLineIndex = gTerminalState.MenuTopRowIndex + gTerminalState.MenuHeight
	moveCursorToLeft()
	moveCursorToPreviousNthLines(gTerminalState.MenuHeight - 1)
	gTerminalState.SelectedLineIndex = gTerminalState.SelectedLineIndex - (gTerminalState.MenuHeight - 1)
}

func dismissMenuLevel1() {
	gTerminalState.CurrentMenuLevel = gTerminalState.CurrentMenuLevel - 1
	jumpCursorToCertainLine(gTerminalState.MenuTopRowIndex)
	//selectedLineIndexBeforeDismiss := gTerminalState.SelectedLineIndex
	// 遍历并清空屏幕
	for i := 0; i < gTerminalState.MenuHeight; i++ {
		clearCurrentLine()
		fmt.Print("\r")
		if i < gTerminalState.MenuHeight-1 {
			moveCursorToNextNthLines(1)
		}
	}
	// 光标停留在menu最后一行，然后移动到menu最上面一行
	moveCursorToPreviousNthLines(gTerminalState.MenuHeight - 1)
	for i := 0; i < gTerminalState.MenuHeight; i++ {
		clearCurrentLine()
		if gTerminalState.SelectedGroupIndex*gTerminalState.SwitchScreenLines+gTerminalState.MenuTopRowIndex+i < len(gSearchData.FileDataArr) {
			printCurrentLineWithUnselectedDisplayName(gTerminalState.SelectedGroupIndex, gTerminalState.MenuTopRowIndex+i)
		}
		fmt.Print("\r")
		if i < gTerminalState.MenuHeight-1 {
			moveCursorToNextNthLines(1)
		}
	}
	// 光标停留在menu最后一行，然后移动到原来的位置
	gTerminalState.SelectedLineIndex = gTerminalState.MenuTopRowIndex + gTerminalState.MenuHeight - 1
	jumpCursorToCertainLine(gTerminalState.MenuLevelOCursorIndex)
	gTerminalState.SelectedLineIndex = gTerminalState.MenuLevelOCursorIndex
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
