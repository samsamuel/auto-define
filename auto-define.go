package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/tidwall/gjson"
)

func main() {
	argLength := len(os.Args[1:])
	if argLength == 0 {
		fmt.Println("Hey! You forgot to give me a file to read.")
		fmt.Println("Should look like this:")
		fmt.Println("/path/to/auto-define /path/to/your-terms.txt")
		os.Exit(0)
	}
	filePath := os.Args[1]
	//fmt.Println(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		term := scanner.Text()
		result := define(term)
		fmt.Println(term + " - " + result)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	input := bufio.NewScanner(os.Stdin)
	fmt.Println("\nPress enter to close...")
	input.Scan()
	fmt.Println(input.Text())
}

func define(x string) (a string) {
	clean_x := strings.ReplaceAll(x, " ", "_")
	url := "https://en.wikipedia.org/w/api.php?action=query&prop=extracts&exsentences=1&exlimit=1&titles=" + clean_x + "&explaintext=1&formatversion=2&format=json"

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	robots, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	json := string(robots[:])

	value := gjson.Get(json, "query.pages.0.extract")

	return value.String()
}
