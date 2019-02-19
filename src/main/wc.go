package main

import (
	"fmt"
	"mapreduce"
	"os"
	"strconv"
	"strings"
	"unicode"
)

//
// The map function is called once for each file of input. The first
// argument is the name of the input file, and the second is the
// file's complete contents. You should ignore the input file name,
// and look only at the contents argument. The return value is a slice
// of key/value pairs.
//
func mapF(filename string, contents string) []mapreduce.KeyValue {
	f := func(c rune) bool {
		return !unicode.IsLetter(c)
	}

	wordsDict := make(map[string]int)

	fields := strings.FieldsFunc(contents, f)
	for _, field := range fields {
		wordsDict[field] += 1
	}

	results := make([]mapreduce.KeyValue, 0)
	for k, v := range wordsDict {
		if len(k) == 0 {
			continue
		}
		kv := mapreduce.KeyValue{Key: k, Value: strconv.Itoa(v)}
		results = append(results, kv)
	}

	return results
}

//
// The reduce function is called once for each key generated by the
// map tasks, with a list of all the values created for that key by
// any map task.
//
func reduceF(key string, values []string) string {
	var sum int = 0

	for _, v := range values {
		fmt.Printf("key is %s, value is %s\n", key, v)

		if amount, err := strconv.Atoi(v); err != nil {
			panic(v)
		} else {
			sum += amount
		}
	}

	return strconv.Itoa(sum)
}

// Can be run in 3 ways:
// 1) Sequential (e.g., go run wc.go master sequential x1.txt .. xN.txt)
// 2) Master (e.g., go run wc.go master localhost:7777 x1.txt .. xN.txt)
// 3) Worker (e.g., go run wc.go worker localhost:7777 localhost:7778 &)
func main() {
	if len(os.Args) < 4 {
		fmt.Printf("%s: see usage comments in file\n", os.Args[0])
	} else if os.Args[1] == "master" {
		var mr *mapreduce.Master
		if os.Args[2] == "sequential" {
			mr = mapreduce.Sequential("wcseq", os.Args[3:], 3, mapF, reduceF)
		} else {
			mr = mapreduce.Distributed("wcseq", os.Args[3:], 3, os.Args[2])
		}
		mr.Wait()
	} else {
		mapreduce.RunWorker(os.Args[2], os.Args[3], mapF, reduceF, 100, nil)
	}
}
