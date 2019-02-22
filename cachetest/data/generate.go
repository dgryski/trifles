package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/pingcap/go-ycsb/pkg/generator"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("error: invalid command")
		fmt.Println("Usage: generate.go <file> <maxKeys> <workloadSize>")
		return
	}

	fileName := os.Args[1]
	maxKeys, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		panic(err)
	}
	workloadSize, err := strconv.ParseInt(os.Args[3], 10, 64)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := genZipfWorkload(file, maxKeys, workloadSize); err != nil {
		panic(err)
	}
}

func genZipfWorkload(file *os.File, maxKeys, workloadSize int64) error {
	z := generator.NewScrambledZipfian(0, maxKeys, generator.ZipfianConstant)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := int64(0); i < workloadSize; i++ {
		if _, err := fmt.Fprintln(file, z.Next(r)); err != nil {
			return err
		}
	}

	return nil
}
