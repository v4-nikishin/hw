package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func getValue(oldV string) string {
	if strings.Contains(oldV, "=") {
		return ""
	}
	newV := strings.ReplaceAll(oldV, string("\x00"), string("\x0A"))
	newV = strings.TrimRight(newV, " ")
	return newV
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	fis, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	env := Environment{}
	for _, fs := range fis {
		fileName := fs.Name()
		filePath := dir + "/" + fileName

		f, err := os.Open(filePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer f.Close()

		sfi, err := os.Stat(filePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if !sfi.Mode().IsRegular() {
			continue
		}

		if sfi.Size() == 0 {
			os.Unsetenv(fileName)
			continue
		}

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			v := getValue(scanner.Text())
			_, ok := os.LookupEnv(fileName)
			env[fileName] = EnvValue{Value: v, NeedRemove: ok}
			break
		}
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
			continue
		}
	}
	return env, nil
}
