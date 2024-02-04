package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var gFileDisplayNamesGroup = make([][]string, 0, 2)
var gSelectedGroupIndex = 0
var gFileDisplayNames = make([]string, 0, 16)
var gFileRealNames = make([]string, 0, 16)
var gSelectedLineIndex = 0
var gTerminalColumnNumber = 0
var gTerminalRowNumber = 0

// 一屏幕多少行显示文件
var gSwitchScreenLines = 0

// 所有屏幕最高的高度
var gMaxLineLength = 0

func main() {
	gTerminalColumnNumber, gTerminalRowNumber = getTerminalColumnsAndRows()
	args := os.Args[1:]
	argsLength := len(args)
	if argsLength == 0 {
		gSwitchScreenLines = gTerminalRowNumber - 2
		printCurrentDirFiles(nil)
	} else if (args[0] == "-help" || args[0] == "-h") && argsLength == 1 {
		printInstructions()
		return
	} else {
		var wordSB strings.Builder
		var firstKeywordIndex = 0
		gSwitchScreenLines = gTerminalRowNumber - 2
		if args[0] == "-help" || args[0] == "-h" {
			printInstructions()
			firstKeywordIndex = 1
			gSwitchScreenLines = gTerminalRowNumber - 2 - len(gInstructions)
		}
		for i := firstKeywordIndex; i < argsLength; i++ {
			wordSB.WriteString(args[i])
		}
		pattern := wordSB.String()
		re := regexp.MustCompile("(?i)" + pattern)
		printCurrentDirFiles(re)
	}
	if len(gFileDisplayNames) == 0 {
		fmt.Println("Search results are Empty.")
		return
	}
	clearPreviousNthLines(getSelectedGroupDisplayFileNamesLength())
	printCurrentLineWithSelectedDisplayName()

	err := keyboard.Open()
	if err != nil {
		fmt.Println("Error opening keyboard:", err)
		os.Exit(1)
	}
	defer keyboard.Close()

	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Println("Error reading keyboard input:", err)
			break
		}

		if key == keyboard.KeyArrowUp {
			if gSelectedLineIndex <= 0 {
				continue
			}
			clearCurrentLine()
			printCurrentLineWithUnselectedDisplayName()
			gSelectedLineIndex = gSelectedLineIndex - 1
			if gSelectedLineIndex < getSelectedGroupDisplayFileNamesLength() {
				clearPreviousNthLines(1)
				printCurrentLineWithSelectedDisplayName()
			}
		} else if key == keyboard.KeyArrowDown {
			if gSelectedLineIndex >= getSelectedGroupDisplayFileNamesLength()-1 {
				continue
			}
			clearCurrentLine()
			printCurrentLineWithUnselectedDisplayName()
			gSelectedLineIndex = gSelectedLineIndex + 1
			clearNextLine()
			fmt.Print("\r")
			printCurrentLineWithSelectedDisplayName()
		} else if key == keyboard.KeyArrowRight {
			// 切屏
			if gSelectedGroupIndex >= getGroupLength()-1 {
				continue
			}
			gSelectedGroupIndex = gSelectedGroupIndex + 1
			if gSelectedLineIndex >= getSelectedGroupDisplayFileNamesLength() {
				deltaUp := gSelectedLineIndex - getSelectedGroupDisplayFileNamesLength() + 1
				gSelectedLineIndex = getSelectedGroupDisplayFileNamesLength() - 1
				// 然后定位光标
				moveCursorToPreviousNthLines(deltaUp)
			}
			// 先记一下刚才的selectedLineIndex，清屏后再回到这个位置
			selectedLineIndex := gSelectedLineIndex

			moveCursorToPreviousNthLines(selectedLineIndex)

			for i := 0; i < gMaxLineLength; i++ {
				clearCurrentLine()
				fmt.Print("\r")
				if i < gMaxLineLength-1 {
					moveCursorToNextNthLines(1)
				}
			}
			moveCursorToPreviousNthLines(gMaxLineLength - 1)
			for i := 0; i < getSelectedGroupDisplayFileNamesLength(); i++ {
				clearCurrentLine()
				gSelectedLineIndex = i
				printCurrentLineWithUnselectedDisplayName()
				fmt.Print("\r")
				if gSelectedLineIndex < getSelectedGroupDisplayFileNamesLength()-1 {
					moveCursorToNextNthLines(1)
				}
			}

			moveCursorToPreviousNthLines(gSelectedLineIndex - selectedLineIndex)
			gSelectedLineIndex = selectedLineIndex
			clearCurrentLine()
			printCurrentLineWithSelectedDisplayName()
		} else if key == keyboard.KeyArrowLeft {
			// 切屏
			if gSelectedGroupIndex <= 0 {
				continue
			}
			gSelectedGroupIndex = gSelectedGroupIndex - 1
			// 先记一下刚才的selectedLineIndex，清屏后再回到这个位置
			selectedLineIndex := gSelectedLineIndex

			moveCursorToPreviousNthLines(selectedLineIndex)

			for i := 0; i < gMaxLineLength; i++ {
				clearCurrentLine()
				fmt.Print("\r")
				if i < gMaxLineLength-1 {
					moveCursorToNextNthLines(1)
				}
			}
			moveCursorToPreviousNthLines(gMaxLineLength - 1)
			for i := 0; i < getSelectedGroupDisplayFileNamesLength(); i++ {
				clearCurrentLine()
				gSelectedLineIndex = i
				printCurrentLineWithUnselectedDisplayName()
				fmt.Print("\r")
				if gSelectedLineIndex < getSelectedGroupDisplayFileNamesLength()-1 {
					moveCursorToNextNthLines(1)
				}
			}

			moveCursorToPreviousNthLines(gSelectedLineIndex - selectedLineIndex)
			gSelectedLineIndex = selectedLineIndex
			clearCurrentLine()
			printCurrentLineWithSelectedDisplayName()
		} else if key == keyboard.KeyEsc {
			clearCurrentLine()
			printCurrentLineWithUnselectedDisplayName()
			fmt.Print("\r")
			clearNextNthLine(gMaxLineLength - gSelectedLineIndex)
			break
		} else if key == keyboard.KeyEnter {
			ret := selectCurrentFile()
			if ret == 0 {
				clearNextNthLine(gMaxLineLength - gSelectedLineIndex)
				break
			}
		}
	}
}

func moveCursorToPreviousNthLines(nth int) {
	if nth > 0 {
		fmt.Print("\033[")
		fmt.Print(nth)
		fmt.Print("F")
	}
}

func moveCursorToNextNthLines(nth int) {
	fmt.Print("\033[") // ANSI escape code to move the cursor down one line
	fmt.Print(nth)
	fmt.Print("B")
}

func clearPreviousNthLines(nth int) {
	fmt.Print("\033[")
	fmt.Print(nth)
	fmt.Print("F")       // ANSI escape code to move the cursor up
	fmt.Print("\033[2K") // ANSI escape code to clear the line
}

func clearCurrentLine() {
	fmt.Print("\033[2K") // ANSI escape code to clear the line
	fmt.Print("\r")
}

func clearNextLine() {
	fmt.Print("\033[1B") // ANSI escape code to move the cursor down one line
	fmt.Print("\033[2K") // ANSI escape code to clear the line
}

func clearNextNthLine(nth int) {
	fmt.Print("\033[") // ANSI escape code to move the cursor down one line
	fmt.Print(nth)
	fmt.Print("B")
	fmt.Print("\033[2K") // ANSI escape code to clear the line
}

func printCurrentLineWithUnselectedDisplayName() {
	fmt.Print(getUnselectedFileNameByIndex(gFileDisplayNames, gSelectedGroupIndex, gSelectedLineIndex))
}

func printCurrentLineWithSelectedDisplayName() {
	fmt.Print(getSelectedFileNameByIndex(gFileDisplayNames, gSelectedGroupIndex, gSelectedLineIndex))
	fmt.Print("\r")
}

func printCurrentDirFiles(reg *regexp.Regexp) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	files, err := os.ReadDir(currentDir)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, file := range files {
		fileName := file.Name()
		if !file.IsDir() {
			prepareMatchedFileName(reg, fileName, "[F]")
		} else {
			prepareMatchedFileName(reg, fileName, "[D]")
		}
	}
	prepareMatchedFileNameGroup()
	displayCurrentFileNamesForFirstTime()
}

func displayCurrentFileNamesForFirstTime() {
	// 获取一屏幕
	gSelectedGroupIndex = 0
	gSelectedLineIndex = 0
	currentScreenFileNames := getSelectedGroupDisplayFileNames()
	for _, value := range currentScreenFileNames {
		fmt.Println(value)
	}
}

func getSelectedGroupDisplayFileNamesLength() int {
	if len(gFileDisplayNamesGroup) == 0 {
		return 0
	}
	return len(gFileDisplayNamesGroup[gSelectedGroupIndex])
}

func getGroupLength() int {
	return len(gFileDisplayNamesGroup)
}

func getSelectedGroupDisplayFileNames() []string {
	if len(gFileDisplayNamesGroup) == 0 {
		return nil
	}
	return gFileDisplayNamesGroup[gSelectedGroupIndex]
}

func prepareMatchedFileNameGroup() {
	groupLength := ceil(len(gFileDisplayNames), gSwitchScreenLines)
	for i := 0; i < groupLength; i++ {
		oneScreenContent := make([]string, 0, gSwitchScreenLines)
		for j := 0; j < gSwitchScreenLines; j++ {
			indexOfFileName := i*gSwitchScreenLines + j
			if indexOfFileName < len(gFileDisplayNames) {
				oneScreenContent = append(oneScreenContent, gFileDisplayNames[indexOfFileName])
			} else {
				break
			}
		}
		gFileDisplayNamesGroup = append(gFileDisplayNamesGroup, oneScreenContent)
	}
	if groupLength == 0 {
		gMaxLineLength = 0
	} else if groupLength == 1 {
		gMaxLineLength = len(gFileDisplayNames)
	} else {
		gMaxLineLength = gSwitchScreenLines
	}
}

func prepareMatchedFileName(reg *regexp.Regexp, fileName string, prefix string) {
	if reg == nil {
		str := prefix + "    " + getDisplayFileName(fileName)
		gFileDisplayNames = append(gFileDisplayNames, str)
		gFileRealNames = append(gFileRealNames, fileName)
	} else {
		match := reg.MatchString(fileName)
		if match {
			str := prefix + "    " + getDisplayFileName(fileName)
			gFileDisplayNames = append(gFileDisplayNames, str)
			gFileRealNames = append(gFileRealNames, fileName)
		}
	}
}

func getDisplayFileName(fileName string) string {
	return truncateString(fileName, (gTerminalColumnNumber-4)*9/10)
}

func addStringIntoAString(originalStr string, insertedIndex int, insertedString string) string {
	return originalStr[0:insertedIndex] + insertedString + originalStr[insertedIndex+1:]
}

func getSelectedFileNameByIndex(arr []string, selectedGroupIndex int, selectedLineIndex int) string {
	return addStringIntoAString(arr[selectedGroupIndex*gSwitchScreenLines+selectedLineIndex], 4, ">")
}
func getUnselectedFileNameByIndex(arr []string, selectedGroupIndex int, selectedLineIndex int) string {
	return arr[selectedGroupIndex*gSwitchScreenLines+selectedLineIndex]
}

func selectCurrentFile() int {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return 1
	}

	cmd := exec.Command("open", filepath.Join(currentDir, gFileRealNames[gSelectedGroupIndex*gSwitchScreenLines+gSelectedLineIndex]))

	// 获取命令的输出
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return 2
	}
	return 0
}
