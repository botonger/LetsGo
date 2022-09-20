package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func listFile(path string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(path)
	checkErr(err)
	return files, nil
}

func main() {

	var originDir string
	var targetDir string

	flag.StringVar(&originDir, "origin", "", "origin directory")
	flag.StringVar(&targetDir, "target", "", "target directory")
	flag.Parse()

	files, err := listFile(originDir)
	checkErr(err)

	for _, file := range files {
		readFile, err := os.Open(originDir + "/" + file.Name())
		checkErr(err)

		writeFile, err := os.Create(targetDir + "/" + file.Name())
		checkErr(err)

		s := bufio.NewScanner(readFile)
		isResponse, err := regexp.MatchString(".*Response.java", file.Name())

		if isResponse {
			for s.Scan() {
				_, err = fmt.Fprintln(writeFile, s.Text())
				checkErr(err)
			}
		} else {
			for s.Scan() {
				isFieldMatched, err := regexp.MatchString(".*private.*null;", s.Text())
				checkErr(err)
				isClassMatched, err := regexp.MatchString("public class ", s.Text())
				checkErr(err)

				var line string

				if isClassMatched {
					line = "@ApiModel(description= \"\")\n" + s.Text() + "\n"
				}
				if isFieldMatched {
					line = "    @ApiModelProperty(value = \"\", notes = \"\", example = \"\")\n" + s.Text()
				}
				if !isFieldMatched && !isClassMatched {
					line = s.Text()
				}
				_, err = fmt.Fprintln(writeFile, line)
				checkErr(err)
			}
		}
	}
}
