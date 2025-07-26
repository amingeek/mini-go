package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func GetList(r io.Reader) []string {
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	words := strings.Fields(scanner.Text())
	return words[1:]
}

func GetSeason(r io.Reader) string {
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	words := strings.Fields(scanner.Text())
	return words[0]
}

func PrintList(shirt []string, pants []string, jacket []string, coat []string, caps []string, season string) {
	for _, i := range shirt {
		for _, j := range pants {
			if season == "SUMMER" || season == "SPRING" || season == "FALL" {
				for _, k := range caps {
					fmt.Printf("SHIRT: %s PANTS: %s CAP: %s\n", i, j, k)
				}
			}
			if season == "SPRING" || season == "FALL" || season == "WINTER" {
				for _, k := range caps {

					for _, x := range coat {
						if season == "FALL" && x == "yellow" || season == "FALL" && x == "orange" {
							continue
						}
						fmt.Printf("COAT: %s SHIRT: %s PANTS: %s CAP: %s\n", x, i, j, k)
					}

				}
				if season == "WINTER" {
					for _, x := range jacket {
						fmt.Printf("SHIRT: %s PANTS: %s JACKET: %s\n", i, j, x)
					}
				}
				for _, x := range coat {
					fmt.Printf("COAT: %s SHIRT: %s PANTS: %s\n", x, i, j)
				}
				if season != "WINTER" {
					fmt.Printf("SHIRT: %s PANTS: %s\n", i, j)
				}
			}

		}
	}
}

func main() {
	coat := GetList(os.Stdin)
	shirt := GetList(os.Stdin)
	pants := GetList(os.Stdin)
	caps := GetList(os.Stdin)
	jacket := GetList(os.Stdin)
	season := GetSeason(os.Stdin)

	PrintList(shirt, pants, jacket, coat, caps, season)
}
