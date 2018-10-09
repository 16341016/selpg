package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"

	flag "github.com/spf13/pflag"
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
	flag.IntVarP(&args.s, "start", "s", 0, "start")
	flag.IntVarP(&args.e, "end", "e", 0, "end")
	flag.IntVarP(&args.l, "line", "l", -1, "line")
	flag.BoolVarP(&args.f, "final", "f", false, "final")
	flag.StringVarP(&args.d, "inputFile", "d", "", "input file")
	flag.Parse()
	otherArgs := flag.Args()
	if len(otherArgs) > 0 {
		args.inputFile = otherArgs[0]
	} else {
		args.inputFile = ""
	}
}

func checkArgs(args *Args) {
	if args.s == 0 || args.e == 0 {
		os.Stderr.Write([]byte("Please input -s and -e\n"))
		os.Exit(0)
	}
	if args.s > args.e {
		os.Stderr.Write([]byte("Invalid input about -s and -e\n"))
		os.Exit(0)
	}
	if args.f && args.l != -1 {
		os.Stderr.Write([]byte("Please choose either -f or -l\n"))
		os.Exit(0)
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
			if args.l == -1 {
				args.l = 72
			}
			readByL(args, reader, writer)
		}
	} else {
		cmd := exec.Command("./" + args.d)
		writer, err := cmd.StdinPipe()
		if err != nil {
			fmt.Println("Error", err)
			os.Exit(1)
		}
		if err := cmd.Start(); err != nil {
			fmt.Println("Error", err)
			os.Exit(1)
		}
		if args.f {
			readByFWithD(args, reader, writer)
		} else {
			if args.l == -1 {
				args.l = 72
			}
			readByLWithD(args, reader, writer)
		}
		writer.Close()
		if err := cmd.Wait(); err != nil {
			fmt.Println("Error")
			os.Exit(1)
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
					os.Exit(1)
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
				os.Exit(1)
			}
			if pageCount >= args.s {
				writer.WriteByte(char)
				writer.Flush()
			}
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
					os.Exit(1)
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
				os.Exit(1)
			}
			if pageCount >= args.s {
				writer.Write([]byte{char})
			}
		}
	}
}

func main() {
	args := new(Args)
	getArgs(args)
	checkArgs(args)
	execution(args)
}
