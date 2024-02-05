package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var gCommandHelp = map[string]int{
	"-help": 0,
	"-h":    0,
}

var gCommandRecursive = map[string]int{
	"-recursive": 0,
	"-r":         0,
}

var gCommands = map[string]int{
	"-help":      0,
	"-h":         0,
	"-recursive": 0,
	"-r":         0,
}

var gInstructions = []string{
	"|------------------------ Instructions ------------------------|",
	"| 1. Press ESC to quit.                                        |",
	"| 2. Press ↑ or ↓ to select a file.                            |",
	"| 3. Press ← or → to switch screen.                            |",
	"| 4. Press Enter to open the selected file.                    |",
	"| 4. Press Space to open the selected file's Directory.        |",
	"| 5. Add -r to recursively walks through current directories.  |",
	"|--------------------------------------------------------------|"}

func printInstructions() {
	for _, value := range gInstructions {
		fmt.Println(value)
	}
}

func containArgHelp(arr []string) bool {
	for _, value := range arr {
		if isArgHelp(value) {
			return true
		}
	}
	return false
}

func containArgRecursive(arr []string) bool {
	for _, value := range arr {
		if isArgRecursive(value) {
			return true
		}
	}
	return false
}

func isArgHelp(s string) bool {
	_, exists := gCommandHelp[s]
	if exists {
		return true
	}
	return false
}

func isArgRecursive(s string) bool {
	_, exists := gCommandRecursive[s]
	if exists {
		return true
	}
	return false
}

func isArgCommand(s string) bool {
	_, exists := gCommands[s]
	if exists {
		return true
	}
	return false
}

func getSearchPattern(arr []string) string {
	filteredArray := make([]string, 0, len(arr))
	for _, value := range arr {
		if isArgCommand(value) {
			continue
		}
		filteredArray = append(filteredArray, value)
	}
	var wordSB strings.Builder
	length := len(filteredArray)
	for index, value := range filteredArray {
		wordSB.WriteString(strings.Replace(value, ".", "\\.", -1))
		if index < length-1 {
			wordSB.WriteString(".*")
		}
	}
	return wordSB.String()
}

func getTerminalColumns() (int, error) {
	var cols int
	cmd := exec.Command("tput", "cols")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err == nil {
		cols, err = strconv.Atoi(strings.TrimSpace(string(out)))
	}
	return cols, err
}

func getTerminalRows() (int, error) {
	var cols int
	cmd := exec.Command("tput", "lines")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err == nil {
		cols, err = strconv.Atoi(strings.TrimSpace(string(out)))
	}
	return cols, err
}

func getTerminalColumnsAndRows() (int, int) {
	cols, err := getTerminalColumns()

	if err != nil {
		cols = 0
	}
	rows, err := getTerminalRows()
	if err != nil {
		rows = 0
	}
	if cols == 0 {
		cols = 80
	}
	if rows == 0 {
		rows = 24
	}
	return cols, rows
}

func truncateString(input string, maxLength int) string {
	if len(input) <= maxLength {
		return input
	}
	// 中间四个....
	halfMaxLength := maxLength/2 - 2

	runeSlice := []rune(input)
	leftRuneByteLength := 0
	rightRuneByteLength := 0
	var leftRuneSlice []rune
	var rightRuneSlice []rune
	for i, r := range runeSlice {
		byteLength := utf8.RuneLen(r)
		leftRuneByteLength = leftRuneByteLength + byteLength
		if leftRuneByteLength >= halfMaxLength {
			leftRuneSlice = runeSlice[0:i]
			break
		}
	}
	for i := len(runeSlice) - 1; i >= 0; i-- {
		r := runeSlice[i]
		byteLength := utf8.RuneLen(r)
		rightRuneByteLength = rightRuneByteLength + byteLength
		if rightRuneByteLength >= halfMaxLength {
			rightRuneSlice = runeSlice[i+1:]
			break
		}
	}
	return string(leftRuneSlice) + "...." + string(rightRuneSlice)
}

func ceil(numerator int, denominator int) int {
	result := numerator / denominator
	remainder := numerator % denominator
	if remainder > 0 {
		result++
	}
	return result
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
	if gContainRecursive {
		err = filepath.Walk(currentDir, func(path string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
				return nil
			}
			filePath := path
			if len(path) > len(currentDir)+1 {
				filePath = path[len(currentDir)+1:]
			}
			if !fileInfo.IsDir() {
				prepareMatchedFileName(reg, filePath, fileInfo.Name(), "[F]")
			} else {
				prepareMatchedFileName(reg, filePath, fileInfo.Name(), "[D]")
			}
			return nil
		})
	} else {
		files, err := os.ReadDir(currentDir)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		for _, file := range files {
			fileName := file.Name()
			if !file.IsDir() {
				prepareMatchedFileName(reg, fileName, fileName, "[F]")
			} else {
				prepareMatchedFileName(reg, fileName, fileName, "[D]")
			}
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

func prepareMatchedFileName(reg *regexp.Regexp, filePath string, fileName string, prefix string) {
	if reg == nil {
		str := prefix + "    " + getDisplayFileName(filePath)
		gFileDisplayNames = append(gFileDisplayNames, str)
		gFileRealNames = append(gFileRealNames, filePath)
	} else {
		match := reg.MatchString(fileName)
		if match {
			str := prefix + "    " + getDisplayFileName(filePath)
			gFileDisplayNames = append(gFileDisplayNames, str)
			gFileRealNames = append(gFileRealNames, filePath)
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
		fmt.Println("Error occurred when executing command:", err)
		return 2
	}
	return 0
}

func selectCurrentFilesParentDir() int {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return 1
	}
	cmd := exec.Command("open",
		filepath.Dir(filepath.Join(currentDir, gFileRealNames[gSelectedGroupIndex*gSwitchScreenLines+gSelectedLineIndex])))
	// 获取命令的输出
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error occurred when executing command:", err)
		return 2
	}
	return 0
}
