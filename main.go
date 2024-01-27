package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	args := os.Args[1:]
	argsLength := len(args)
	if argsLength == 0 {
		printCurrentDirFiles(nil)
		return
	}

	var wordSB strings.Builder
	for i := 0; i < argsLength; i++ {
		wordSB.WriteString(args[i])
	}

	// 创建正则表达式对象
	pattern := wordSB.String()
	re := regexp.MustCompile(pattern)
	printCurrentDirFiles(re)
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
		fmt.Println("[F]	", fileName)
	} else {
		match := reg.MatchString(fileName)
		if match {
			fmt.Println("[F]	", fileName)
		}
	}
}
