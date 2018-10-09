package main

import (
	"bufio"
	"io"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	file, err := os.OpenFile("./testD.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		os.Stderr.WriteString("Open Failed")
	}
	writer := bufio.NewWriter(file)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			os.Stderr.Write([]byte("Read error\n"))
			os.Exit(1)
		}
		writer.Write(line)
		writer.Flush()
	}
}
