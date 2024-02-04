package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unicode/utf8"
)

var gInstructions = []string{
	"|---------------- Instructions ----------------|",
	"| 1. Press ESC to quit.                        |",
	"| 2. Press ↑ or ↓ to select a file.            |",
	"| 3. Press ← or → to switch screen.            |",
	"| 4. Press Enter to open the selected file.    |",
	"|----------------------------------------------|"}

func printInstructions() {
	for _, value := range gInstructions {
		fmt.Println(value)
	}
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

//func getCursorLine() (string, error) {
//	cmd := exec.Command("tput", "cup")
//	output, err := cmd.CombinedOutput()
//	if err != nil {
//		return "", err
//	}
//	fmt.Println("--->" + string(output) + "<---")
//	// 解析命令输出，获取行号
//	values := strings.Fields(string(output))
//	if len(values) != 2 {
//		return "", fmt.Errorf("unexpected output format")
//	}
//
//	return values[0], nil
//}
