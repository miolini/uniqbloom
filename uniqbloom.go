package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	
	"github.com/willf/bloom"
)

const (
	newLine = '\n'
)

var (
	flBloomN = flag.Uint("n", 1e6, "bloom filter n arg")
	flBloomE = flag.Float64("e", 1e-5, "bloom filter error factor")
	flQuiet  = flag.Bool("q", true, "quiet mode")
)

func main() {
	var (
		line []byte
		err  error
	)
	if !(*flQuiet) {
		fmt.Fprintf(os.Stderr, "Filter repeated lines by Artem Andreenko <mio@volmy.com>\n")
	}
	flag.Parse()
	filter := bloom.NewWithEstimates(uint(*flBloomN), *flBloomE)
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for {
		line, err = reader.ReadBytes(newLine)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "stdin read error: %s\n", err)
		}
		if filter.TestAndAdd(line) {
			continue
		}
		_, err = writer.Write(line)
		if err != nil {
			if err == io.EOF {
				fmt.Fprintf(os.Stderr, "waring: stdout closed before stdin\n")
				break
			}
			fmt.Fprintf(os.Stderr, "write to stdout error: %s\n", err)
			break
		}
	}
}
