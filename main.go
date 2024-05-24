package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const LOGO string = " _                           _ __     __    _ _     _       _             \n| |__   ___ _ __ _   _ _ __ | |\\ \\   / /_ _| (_) __| | __ _| |_ ___  _ __ \n| '_ \\ / __| '__| | | | '_ \\| __\\ \\ / / _` | | |/ _` |/ _` | __/ _ \\| '__|\n| |_) | (__| |  | |_| | |_) | |_ \\ V / (_| | | | (_| | (_| | || (_) | |   \n|_.__/ \\___|_|   \\__, | .__/ \\__| \\_/ \\__,_|_|_|\\__,_|\\__,_|\\__\\___/|_|   \n                 |___/|_|                                                 \n"

type PasswordRecord struct {
	Email    string
	HashPass string
}

var mutex sync.Mutex

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
	threadNum := flag.Int("thread", 4, "运行线程数")
	flag.Parse()

	if *inputFile == "" {
		fmt.Printf(LOGO)
		fmt.Printf("请提供输入文件路径")
		return
	}

	if *passwordToVerify == "" {
		fmt.Printf(LOGO)
		fmt.Printf("请提供需要验证的密码")
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
	err = file.Close()
	if err != nil {
		return
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
	defer func(resultFile *os.File) {
		err := resultFile.Close()
		if err != nil {

		}
	}(resultFile)

	// 创建一个带缓冲的通道
	passwordRecordChan := make(chan PasswordRecord, 20)
	numWorkers := *threadNum
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	fmt.Println("Worker count:", numWorkers)

	// 启动工作goroutine
	for i := 0; i < numWorkers; i++ {

		go func() {
			defer wg.Done()
			for record := range passwordRecordChan {
				if verify(*passwordToVerify, record.HashPass) {
					mutex.Lock()
					fmt.Printf("Verification succeeded: %s %s --> %s\n", record.Email, *passwordToVerify, record.HashPass)
					_, err := resultFile.WriteString(fmt.Sprintf("%s,%s,%s\n", record.Email, *passwordToVerify, record.HashPass))
					if err != nil {
						fmt.Println(err.Error())
						continue
					}
					mutex.Unlock()
				}
			}
		}()
	}
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		fmt.Println("Program exiting...")
		close(passwordRecordChan)
		err := resultFile.Close()
		if err != nil {
			return
		}
		os.Exit(0)
	}()
	// 将密码发送到通道
	for _, record := range passwordList {
		passwordRecordChan <- record
	}

	// 关闭通道
	close(passwordRecordChan)

	// 等待所有工作goroutine完成
	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("Verifying %d passwords took %.1f seconds\n", len(passwordList), elapsed.Seconds())
}
