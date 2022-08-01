package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
	1.启动一个协程，writeDataToFile，随机生成1000个数据
	2.当writeDataToFile写入到文件后，让sort协程从文件中读取这1000个数据完成排序，并输出到另一个文件
	3.扩展：开启10个writeDataToFile存放到十个文件中，
			10个sort协程从十个文件中读取并完成排序，重新写入到十个结果文件中
*/
func main() {
	writeToFileChan := make(chan int, 10)
	exitChan := make(chan bool, 10)
	sortChan := make(chan bool, 10)
	var filePath string
	for i := 0; i < 10; i++ {
		filePath = fmt.Sprintf("E:\\TEST\\%d.txt", i)
		go writeDataToFile(writeToFileChan, filePath, i, exitChan)
	}
	// writeToFileChan多会关？
	go func() {
		// 通过exiChant判断  无需关闭，等着就行
		for i := 0; i < 10; i++ {
			<-exitChan
		}
		close(writeToFileChan)
		fmt.Println("写入完毕")
	}()
	// 读取writeToFileChan，close才有意义！！！
	for i := range writeToFileChan {
		filePath = fmt.Sprintf("E:\\TEST\\%d.txt", i)
		go sortFromNum(sortChan, filePath, i)
	}
	// 通过sortChan阻塞判断十个文件是否排序完毕
	for i := 0; i < 10; i++ {
		<-sortChan // 阻塞
	}
	fmt.Println("排序完毕")
	fmt.Println("主线程退出")

}

func writeDataToFile(writeToFileChan chan int, filePath string, i int, exitChan chan bool) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("open file err=%v \n", err)
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for i := 1; i <= 30; i++ {
		if i%10 == 0 {
			writer.WriteString(fmt.Sprintf("%d\n", rand.Intn(10000)))
		} else {
			writer.WriteString(fmt.Sprintf("%d,", rand.Intn(10000)))
		}
	}
	writer.Flush()
	writeToFileChan <- i
	exitChan <- true
}

func sortFromNum(sortChan chan bool, filePath string, i int) {
	// 从writeToFileChan取出一个，如果取完则直接退出
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("open file err=", err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	// 读取文件内容
	var numSlice []int
	for {
		//读到一个换行符就结束
		readString, err := reader.ReadString('\n')
		// io.EOF表示文件末尾
		if err == io.EOF {
			break
		}
		strArr := strings.Split(strings.Trim(readString, "\n"), ",")
		tempSlice := make([]int, len(strArr), 100)
		for i, v := range strArr {
			tempSlice[i], err = strconv.Atoi(v)
		}
		// 追加
		numSlice = append(numSlice, tempSlice...)
	}
	// 排序
	sort.Ints(numSlice)
	// 写入文件
	resFilePath := fmt.Sprintf("E:\\TEST\\res%d.txt", i)
	file, err = os.OpenFile(resFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("open file err=%v \n", err)
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	str := strings.Trim(fmt.Sprint(numSlice), "[]")
	fmt.Println("排序后：", str)
	writer.WriteString(str)
	writer.Flush()
	// 管道通知
	sortChan <- true
}
