package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

var gCommandHelp = map[string]interface{}{
	"-help": struct{}{},
	"-h":    struct{}{},
}

var gCommandPlain = map[string]interface{}{
	"-plain": struct{}{},
	"-p":     struct{}{},
}

var gCommandTime = map[string]interface{}{
	"-time": struct{}{},
	"-t":    struct{}{},
}

var gCommandCount = map[string]interface{}{
	"-count": struct{}{},
	"-c":     struct{}{},
}

var gCommands = map[string]interface{}{
	"-help":  struct{}{},
	"-h":     struct{}{},
	"-plain": struct{}{},
	"-p":     struct{}{},
	"-time":  struct{}{},
	"-t":     struct{}{},
}

var gInstructions = []string{
	"|------------------------ Instructions ------------------------|",
	"| 1. Press ESC to quit.                                        |",
	"| 2. Press ↑ or ↓ to select a file.                            |",
	"| 3. Press ← or → to switch screen.                            |",
	"| 4. Press Enter to open the selected file.                    |",
	"| 5. Press Space to open the selected file's Directory.        |",
	"| 6. Add -p to search the plain first layer Directory.         |",
	"| 7. Add -t to display files in update time order.             |",
	"| 8. Add -cx to search the first x files and stop.             |",
	"|--------------------------------------------------------------|"}

var gMenu = []string{
	" +------------------------+ ",
	" |    Info            [I] | ",
	" |    Rename          [R] | ",
	" |    Delete          [D] | ",
	" |    Parent Folder   [P] | ",
	" |    Close           [C] | ",
	" +------------------------+ ",
}

func printHelpInstructions() {
	for _, value := range gInstructions {
		fmt.Println(value)
	}
}

func isArgHelp(s string) bool {
	_, exists := gCommandHelp[s]
	if exists {
		return true
	}
	return false
}

func isArgPlain(s string) bool {
	_, exists := gCommandPlain[s]
	if exists {
		return true
	}
	return false
}

func isArgTime(s string) bool {
	_, exists := gCommandTime[s]
	if exists {
		return true
	}
	return false
}

func getSearchCount(s string) int {
	count := 0
	for key, _ := range gCommandCount {
		if strings.HasPrefix(s, key) {
			count, _ = strconv.Atoi(s[len(key):])
			if count > 0 {
				return count
			}
		}
	}
	return count
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

func getTerminalColumnsAndRows() {
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
	gTerminalState.TerminalColumnNumber = cols
	gTerminalState.TerminalRowNumber = rows
}

func getCommandState() {
	args := os.Args[1:]
	argsLength := len(args)
	if argsLength == 0 {
		return
	}
	argsForSearchPattern := make([]string, 0, argsLength)
	for _, arg := range args {
		_, exists := gCommands[arg]
		if exists {
			switch {
			case isArgPlain(arg):
				gCommandState.Plain = true
			case isArgHelp(arg):
				gCommandState.Help = true
			case isArgTime(arg):
				gCommandState.Time = true
			}
		} else {
			searchCount := getSearchCount(arg)
			if searchCount > 0 {
				gCommandState.Count.CountSwitch = true
				gCommandState.Count.CountNumber = searchCount
			} else {
				argsForSearchPattern = append(argsForSearchPattern, arg)
			}
		}
	}
	var patternBuilder strings.Builder
	length := len(argsForSearchPattern)
	for index, value := range argsForSearchPattern {
		patternBuilder.WriteString(strings.Replace(value, ".", "\\.", -1))
		if index < length-1 {
			patternBuilder.WriteString(".*")
		}
	}
	gCommandState.SearchPattern = patternBuilder.String()
}

func getScreenLineNumber() {
	if !gCommandState.Help {
		// 命令不包含help
		gTerminalState.SwitchScreenLines = gTerminalState.TerminalRowNumber - 2
	} else {
		// 命令包含help
		gTerminalState.SwitchScreenLines = gTerminalState.TerminalRowNumber - 2 - len(gInstructions)
	}
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

func jumpCursorToCertainLine(destLineIndex int) {
	if destLineIndex < gTerminalState.SelectedLineIndex {
		moveCursorToPreviousNthLines(gTerminalState.SelectedLineIndex - destLineIndex)
	} else {
		moveCursorToNextNthLines(destLineIndex - gTerminalState.SelectedLineIndex)
	}
}

func moveCursorToNextNthLines(nth int) {
	fmt.Print("\033[") // ANSI escape code to move the cursor down one line
	fmt.Print(nth)
	fmt.Print("B")
}

func clearPreviousNthLine(nth int) {
	fmt.Print("\033[")
	fmt.Print(nth)
	fmt.Print("F")       // ANSI escape code to move the cursor up
	fmt.Print("\033[2K") // ANSI escape code to clear the line
}

func clearCurrentLine() {
	fmt.Print("\033[2K") // ANSI escape code to clear the line
	fmt.Print("\r")
}
func moveCursorToColumnIndex(columnIndex int) {
	fmt.Print("\r")
	fmt.Printf("\033[%dC", columnIndex)
}

func moveCursorToLeft() {
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

func printContentLineWithUnselectedDisplayName(selectedGroupIndex int, selectedLineIndex int) {
	fmt.Print(getUnselectedDisplayFileNameByIndex(gSearchData.FileDataArr, selectedGroupIndex, selectedLineIndex))
}

func printSelectedMenuLevel1Line(selectedMenuIndex int) {
	selectedMenuLine := addStringIntoAString(gMenu[selectedMenuIndex], 3, ">>")
	fmt.Print(selectedMenuLine)
}

func printUnselectedMenuLevel1Line(selectedMenuIndex int) {
	fmt.Print(gMenu[selectedMenuIndex])
}

func printContentLineWithSelectedDisplayName(selectedGroupIndex int, selectedLineIndex int) {
	fmt.Print(getSelectedDisplayFileNameByIndex(gSearchData.FileDataArr, selectedGroupIndex, selectedLineIndex))
	fmt.Print("\r")
}

func onlyPrintHelpInstructions() bool {
	if len(gCommandState.SearchPattern) == 0 &&
		gCommandState.Help &&
		!gCommandState.Plain &&
		!gCommandState.Time {
		return true
	}
	return false
}

func printCurrentDirFiles(reg *regexp.Regexp) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if gCommandState.Plain {
		// 只搜索第一层
		files, err := os.ReadDir(currentDir)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		for _, file := range files {
			fileName := file.Name()
			fileInfo, err := file.Info()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			prepareMatchedFileInfo(reg, fileName, &fileInfo)
			if gCommandState.Count.CountSwitch {
				if len(gSearchData.FileDataArr) >= gCommandState.Count.CountNumber {
					break
				}
			}
		}
	} else {
		// 搜索遍历所有层级
		err = filepath.Walk(currentDir, func(path string, fileInfo os.FileInfo, err error) error {
			if path == currentDir {
				return nil
			}
			if err != nil {
				fmt.Println(err)
				return nil
			}
			filePath := path
			if len(path) > len(currentDir)+1 {
				filePath = path[len(currentDir)+1:]
			}
			prepareMatchedFileInfo(reg, filePath, &fileInfo)
			if gCommandState.Count.CountSwitch {
				if len(gSearchData.FileDataArr) >= gCommandState.Count.CountNumber {
					return filepath.SkipDir
				}
			}
			return nil
		})
	}
	// 按照时间排序
	if gCommandState.Time {
		sort.Sort(FileDataSlice(gSearchData.FileDataArr))
	}

	prepareMatchedFileNameGroup()
	displayCurrentFileNamesForFirstTime()
}

func displayCurrentFileNamesForFirstTime() {
	// 获取一屏幕
	gTerminalState.SelectedGroupIndex = 0
	gTerminalState.SelectedLineIndex = 0
	currentScreenFileNames := getSelectedGroupDisplayFileNames(gSearchData.DisplayFileNamesInGroup)
	for _, value := range currentScreenFileNames {
		fmt.Println(value)
	}
}

func getSelectedGroupDisplayFileNamesLength(displayFileNamesInGroup [][]string) int {
	if len(displayFileNamesInGroup) == 0 {
		return 0
	}
	return len(displayFileNamesInGroup[gTerminalState.SelectedGroupIndex])
}

func getGroupLength(displayFileNamesInGroup [][]string) int {
	return len(displayFileNamesInGroup)
}

func getSelectedGroupDisplayFileNames(displayFileNamesInGroup [][]string) []string {
	if len(displayFileNamesInGroup) == 0 {
		return nil
	}
	return displayFileNamesInGroup[gTerminalState.SelectedGroupIndex]
}

func prepareMatchedFileNameGroup() {
	groupLength := ceil(len(gSearchData.FileDataArr), gTerminalState.SwitchScreenLines)
	for i := 0; i < groupLength; i++ {
		oneScreenContent := make([]string, 0, gTerminalState.SwitchScreenLines)
		for j := 0; j < gTerminalState.SwitchScreenLines; j++ {
			indexOfFileName := i*gTerminalState.SwitchScreenLines + j
			if indexOfFileName < len(gSearchData.FileDataArr) {
				oneScreenContent = append(oneScreenContent, gSearchData.FileDataArr[indexOfFileName].DisplayFileName)
			} else {
				break
			}
		}
		gSearchData.DisplayFileNamesInGroup = append(gSearchData.DisplayFileNamesInGroup, oneScreenContent)
	}
	if groupLength == 0 {
		gTerminalState.MaxLineLength = 0
	} else if groupLength == 1 {
		gTerminalState.MaxLineLength = len(gSearchData.FileDataArr)
	} else {
		gTerminalState.MaxLineLength = gTerminalState.SwitchScreenLines
	}
}

func prepareMatchedFileInfo(reg *regexp.Regexp, filePath string, fileInfo *os.FileInfo) {
	prefix := "[F]"
	if (*fileInfo).IsDir() {
		prefix = "[D]"
	}
	if reg == nil {
		displayFileName := prefix + "    " + getDisplayFileName(filePath)
		fileData := FileData{
			DisplayFileName: displayFileName,
			FilePath:        filePath,
			Time:            (*fileInfo).ModTime(),
		}
		gSearchData.FileDataArr = append(gSearchData.FileDataArr, fileData)
	} else {
		match := reg.MatchString((*fileInfo).Name())
		if match {
			displayFileName := prefix + "    " + getDisplayFileName(filePath)
			fileData := FileData{
				DisplayFileName: displayFileName,
				FilePath:        filePath,
				Time:            (*fileInfo).ModTime(),
			}
			gSearchData.FileDataArr = append(gSearchData.FileDataArr, fileData)
		}
	}
}

func getDisplayFileName(fileName string) string {
	return truncateString(fileName, (gTerminalState.TerminalColumnNumber-4)*9/10)
}

func addStringIntoAString(originalStr string, insertedIndex int, insertedString string) string {
	return originalStr[0:insertedIndex] + insertedString + originalStr[insertedIndex+len(insertedString):]
}

func getSelectedDisplayFileNameByIndex(arr []FileData, selectedGroupIndex int, selectedLineIndex int) string {
	return addStringIntoAString(arr[selectedGroupIndex*gTerminalState.SwitchScreenLines+selectedLineIndex].DisplayFileName, 4, ">")
}

func getUnselectedDisplayFileNameByIndex(arr []FileData, selectedGroupIndex int, selectedLineIndex int) string {
	return arr[selectedGroupIndex*gTerminalState.SwitchScreenLines+selectedLineIndex].DisplayFileName
}

func openCurrentFile() int {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return 1
	}
	fileName := filepath.Join(currentDir, gSearchData.FileDataArr[gTerminalState.SelectedGroupIndex*gTerminalState.SwitchScreenLines+gTerminalState.SelectedLineIndex].FilePath)
	cmd := exec.Command("open", fileName)
	// 获取命令的输出
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error occurred when executing command: open", fileName)
		return 2
	}
	return 0
}

func openCurrentFilesParentDir() int {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return 1
	}
	fileName := filepath.Dir(filepath.Join(currentDir, gSearchData.FileDataArr[gTerminalState.SelectedGroupIndex*gTerminalState.SwitchScreenLines+gTerminalState.SelectedLineIndex].FilePath))
	cmd := exec.Command("open", fileName)
	// 获取命令的输出
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error occurred when executing command: open", fileName)
		return 2
	}
	return 0
}
