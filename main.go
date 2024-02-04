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

var gFileDisplayNamesGroup = make([][]string, 0, 16)
var gSelectedGroup = 0

var gFileDisplayNames = make([]string, 0, 16)
var gFileRealNames = make([]string, 0, 16)
var gSelectedIndex = 0
var gTerminalColumnNumber = 0
var gTerminalRowNumber = 0
var gSwitchScreenLine = 0

func main() {
	gTerminalColumnNumber, gTerminalRowNumber = getTerminalColumnsAndRows()
	args := os.Args[1:]
	argsLength := len(args)
	if argsLength == 0 {
		printCurrentDirFiles(nil)
	} else if args[0] == "-help" && argsLength == 1 {
		printInstructions()
		return
	} else {
		var wordSB strings.Builder
		var firstKeywordIndex = 0
		if args[0] == "-help" {
			printInstructions()
			firstKeywordIndex = 1
		}
		for i := firstKeywordIndex; i < argsLength; i++ {
			wordSB.WriteString(args[i])
		}
		pattern := wordSB.String()
		re := regexp.MustCompile("(?i)" + pattern)
		printCurrentDirFiles(re)
	}
	if len(gFileDisplayNames) == 0 {
		fmt.Println("Search result is empty.")
		return
	}
	// 定位到文件列表的第一行，如果终端的输出文字大于一屏幕，则光标只能定位到屏幕边缘，不能向上滚屏
	clearPreviousNthLines(len(gFileDisplayNames))
	gSelectedIndex = 0
	fmt.Print(getSelectedFileNameByIndex(gFileDisplayNames, gSelectedIndex))
	fmt.Print("\r")

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
			if gSelectedIndex <= 0 {
				continue
			}
			clearCurrentLine()
			fmt.Print(getUnselectedFileNameByIndex(gFileDisplayNames, gSelectedIndex))
			gSelectedIndex = gSelectedIndex - 1
			if gSelectedIndex < len(gFileDisplayNames) {
				clearPreviousNthLines(1)
				fmt.Print(getSelectedFileNameByIndex(gFileDisplayNames, gSelectedIndex))
				fmt.Print("\r")
			}
		} else if key == keyboard.KeyArrowDown {
			if gSelectedIndex >= len(gFileDisplayNames)-1 {
				continue
			}
			clearCurrentLine()
			fmt.Print(getUnselectedFileNameByIndex(gFileDisplayNames, gSelectedIndex))
			gSelectedIndex = gSelectedIndex + 1
			clearNextLine()
			fmt.Print("\r")
			fmt.Print(getSelectedFileNameByIndex(gFileDisplayNames, gSelectedIndex))
			fmt.Print("\r")
		} else if key == keyboard.KeyEsc {
			clearNextNthLine(len(gFileDisplayNames) - gSelectedIndex)
			break
		} else if key == keyboard.KeyEnter {
			ret := selectCurrentFile()
			if ret == 0 {
				clearNextNthLine(len(gFileDisplayNames) - gSelectedIndex)
				break
			}
		}
	}
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
			printWithRegAndMatchedFileName(reg, fileName, "[F]")
		} else {
			printWithRegAndMatchedFileName(reg, fileName, "[D]")
		}
	}
}

func printWithRegAndMatchedFileName(reg *regexp.Regexp, fileName string, prefix string) {
	if reg == nil {
		str := prefix + "    " + getDisplayFileName(fileName)
		gFileDisplayNames = append(gFileDisplayNames, str)
		gFileRealNames = append(gFileRealNames, fileName)
		fmt.Println(str)
	} else {
		match := reg.MatchString(fileName)
		if match {
			str := prefix + "    " + getDisplayFileName(fileName)
			gFileDisplayNames = append(gFileDisplayNames, str)
			gFileRealNames = append(gFileRealNames, fileName)
			fmt.Println(str)
		}
	}
}

func getDisplayFileName(fileName string) string {
	return truncateString(fileName, (gTerminalColumnNumber-4)*9/10)
}

func addStringIntoAString(originalStr string, insertedIndex int, insertedString string) string {
	return originalStr[0:insertedIndex] + insertedString + originalStr[insertedIndex+1:]
}

func getSelectedFileNameByIndex(arr []string, selectedIndex int) string {
	return addStringIntoAString(arr[selectedIndex], 4, ">")
}
func getUnselectedFileNameByIndex(arr []string, selectedIndex int) string {
	return arr[selectedIndex]
}

func selectCurrentFile() int {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return 1
	}

	cmd := exec.Command("open", filepath.Join(currentDir, gFileRealNames[gSelectedIndex]))

	// 获取命令的输出
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return 2
	}
	return 0
}
