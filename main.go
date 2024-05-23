package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const LOGO string = " _                           _ __     __    _ _     _       _             \n| |__   ___ _ __ _   _ _ __ | |\\ \\   / /_ _| (_) __| | __ _| |_ ___  _ __ \n| '_ \\ / __| '__| | | | '_ \\| __\\ \\ / / _` | | |/ _` |/ _` | __/ _ \\| '__|\n| |_) | (__| |  | |_| | |_) | |_ \\ V / (_| | | | (_| | (_| | || (_) | |   \n|_.__/ \\___|_|   \\__, | .__/ \\__| \\_/ \\__,_|_|_|\\__,_|\\__,_|\\__\\___/|_|   \n                 |___/|_|                                                 \n"

type PasswordRecord struct {
	Email    string
	HashPass string
}

func verify(password, hashPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password))
	return err == nil
}

func main() {
	// 从命令行获取输入文件名
	inputFile := flag.String("input", "", "输入待验证的密码文件路径")
	// 从命令行获取需要验证的密码
	passwordToVerify := flag.String("password", "", "需要验证的密码")
	verifyNum := flag.Int("count", -1, "验证限制数量,默认为-1不限制")
	outputFile := flag.String("output", "", "输出结果文件的路径")
	flag.Parse()

	if *inputFile == "" {
		fmt.Println(LOGO)
		fmt.Println("请提供输入文件路径")
		return
	}

	if *passwordToVerify == "" {
		fmt.Println(LOGO)
		fmt.Println("请提供需要验证的密码")
		return
	}

	// 从命令行获取输出文件名(可选)

	if *outputFile == "" {
		*outputFile = filepath.Base(strings.Replace(*inputFile, ".", "_", -1)) + "_result.txt"
	}
	fmt.Printf("Output file name: %s\n", *outputFile)

	passwordList := make([]PasswordRecord, 0)

	// 读取密码文件
	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if *verifyNum >= 0 && len(passwordList) >= *verifyNum {
			break
		}
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) == 2 {
			email, hashPass := parts[0], parts[1]
			passwordList = append(passwordList, PasswordRecord{Email: email, HashPass: hashPass})
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Printf("Loaded %d password records\n", len(passwordList))

	start := time.Now()

	// 创建结果文件
	resultFile, err := os.Create(*outputFile)
	if err != nil {
		fmt.Println("Error creating result file:", err)
		return
	}
	defer resultFile.Close()

	for _, record := range passwordList {
		if verify(*passwordToVerify, record.HashPass) {
			fmt.Printf("Verification succeeded: %s %s --> %s\n", record.Email, *passwordToVerify, record.HashPass)
			_, err := resultFile.WriteString(fmt.Sprintf("%s,%s,%s\n", record.Email, *passwordToVerify, record.HashPass))
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Verifying %d passwords took %.1f seconds\n", len(passwordList), elapsed.Seconds())
}
