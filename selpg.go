package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Args struct {
	s         int
	e         int
	l         int
	f         bool
	d         string
	inputFile string
}

func getArgs(args *Args) {
	flag.IntVar(&args.s, "s", 0, "start")
	flag.IntVar(&args.e, "e", 0, "end")
	flag.IntVar(&args.l, "l", 72, "line")
	flag.BoolVar(&args.f, "f", false, "final")
	flag.StringVar(&args.d, "d", "", "input file")
	flag.Parse()
	otherArgs := flag.Args()
	if len(otherArgs) > 0 {
		args.inputFile = otherArgs[0]
	} else {
		args.inputFile = ""
	}
}

func getReader(args *Args) *bufio.Reader {
	var reader *bufio.Reader
	if args.inputFile == "" {
		reader = bufio.NewReader(os.Stdin)
	} else {
		file, err := os.Open("./" + args.inputFile)
		if err != nil {
			os.Stderr.Write([]byte("File does not exist\n"))
			os.Exit(1)
		}
		reader = bufio.NewReader(file)
	}
	return reader
}

func execution(args *Args) {
	//get reader
	var reader *bufio.Reader
	reader = getReader(args)

	//get writer
	if args.d == "" {
		writer := bufio.NewWriter(os.Stdout)
		if args.f {
			readByF(args, reader, writer)
		} else {
			readByL(args, reader, writer)
		}
	} else {
		cmd := exec.Command("./" + args.d)
		writer, err := cmd.StdinPipe()
		if err != nil {
			fmt.Println("Error", err)
			os.Exit(2)
		}
		if err := cmd.Start(); err != nil {
			fmt.Println("Error", err)
			os.Exit(2)
		}
		if args.f {
			readByFWithD(args, reader, writer)
		} else {
			readByLWithD(args, reader, writer)
		}
		writer.Close()
		if err := cmd.Wait(); err != nil {
			fmt.Println("Error")
			os.Exit(2)
		}
	}
}

func readByL(args *Args, reader *bufio.Reader, writer *bufio.Writer) {
	for pageCount := 1; pageCount <= args.e; pageCount++ {
		if pageCount < args.s {
			for lineCount := 0; lineCount < args.l; lineCount++ {
				reader.ReadBytes('\n')
			}
		} else {
			for lineCount := 0; lineCount < args.l; lineCount++ {
				line, err := reader.ReadBytes('\n')
				if err != nil {
					if err == io.EOF {
						writer.WriteByte('\n')
						writer.Flush()
						break
					}
					os.Stderr.Write([]byte("Read failed\n"))
					os.Exit(2)
				}
				writer.Write(line)
				writer.Flush()
			}
		}
	}
}

func readByF(args *Args, reader *bufio.Reader, writer *bufio.Writer) {
	for pageCount := 1; pageCount <= args.e; pageCount++ {
		for {
			char, err := reader.ReadByte()
			if char == '\f' {
				break
			}
			if err != nil {
				if err == io.EOF {
					writer.WriteByte('\n')
					writer.Flush()
					break
				}
				os.Stderr.Write([]byte("Read failed\n"))
				os.Exit(2)
			}
			writer.WriteByte(char)
			writer.Flush()
		}
	}
}

func readByLWithD(args *Args, reader *bufio.Reader, writer io.WriteCloser) {
	for pageCount := 1; pageCount <= args.e; pageCount++ {
		if pageCount < args.s {
			for lineCount := 0; lineCount < args.l; lineCount++ {
				reader.ReadBytes('\n')
			}
		} else {
			for lineCount := 0; lineCount < args.l; lineCount++ {
				line, err := reader.ReadBytes('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					os.Stderr.Write([]byte("Read failed\n"))
					os.Exit(2)
				}
				writer.Write(line)
			}
		}
	}
}

func readByFWithD(args *Args, reader *bufio.Reader, writer io.WriteCloser) {
	for pageCount := 1; pageCount <= args.e; pageCount++ {
		for {
			char, err := reader.ReadByte()
			if char == '\f' {
				break
			}
			if err != nil {
				if err == io.EOF {
					break
				}
				os.Stderr.Write([]byte("Read failed\n"))
				os.Exit(2)
			}
			writer.Write([]byte{char})
		}
	}
}

func main() {
	args := new(Args)
	getArgs(args)
	execution(args)
}
