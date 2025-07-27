package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ParseLine(line string) []string {
	words := strings.Fields(line)
	if len(words) == 0 {
		return []string{}
	}
	return words[1:]
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	coat := ParseLine(scanner.Text())

	scanner.Scan()
	shirt := ParseLine(scanner.Text())

	scanner.Scan()
	pants := ParseLine(scanner.Text())

	scanner.Scan()
	caps := ParseLine(scanner.Text())

	scanner.Scan()
	jacket := ParseLine(scanner.Text())

	scanner.Scan()
	seasonLine := strings.Fields(scanner.Text())
	season := ""
	if len(seasonLine) > 0 {
		season = seasonLine[0]
	}

	PrintList(coat, shirt, pants, caps, jacket, season)
}

func PrintList(coat []string, shirt []string, pants []string, caps []string, jacket []string, season string) {
	for _, s := range shirt {
		for _, p := range pants {
			switch season {
			case "SUMMER":
				for _, c := range caps {
					fmt.Printf("SHIRT: %s PANTS: %s CAP: %s\n", s, p, c)
				}
			case "SPRING", "FALL":
				for _, s := range shirt {
					for _, p := range pants {
						// کت‌های مجاز برای FALL بدون زرد و نارنجی
						allowedCoats := coat
						if season == "FALL" {
							allowedCoats = []string{}
							for _, co := range coat {
								if co != "yellow" && co != "orange" {
									allowedCoats = append(allowedCoats, co)
								}
							}
						}

						// حالت 1: بدون کت و با کلاه
						for _, c := range caps {
							fmt.Printf("SHIRT: %s PANTS: %s CAP: %s\n", s, p, c)
						}

						// حالت 2: بدون کت و بدون کلاه
						fmt.Printf("SHIRT: %s PANTS: %s\n", s, p)

						// حالت 3: با کت و بدون کلاه
						for _, co := range allowedCoats {
							fmt.Printf("COAT: %s SHIRT: %s PANTS: %s\n", co, s, p)
						}

						// حالت 4: با کت و کلاه (اگر قانون اجازه می‌دهد)
						// اگر در مسئله نیامده که کت و کلاه نباید همزمان باشند،
						// این حالت را هم اضافه کنیم:
						for _, co := range allowedCoats {
							for _, c := range caps {
								fmt.Printf("COAT: %s SHIRT: %s PANTS: %s CAP: %s\n", co, s, p, c)
							}
						}
					}
				}

			case "WINTER":
				for _, j := range jacket {
					fmt.Printf("SHIRT: %s PANTS: %s JACKET: %s\n", s, p, j)
				}
				for _, co := range coat {
					fmt.Printf("COAT: %s SHIRT: %s PANTS: %s\n", co, s, p)
				}
			}
		}
	}
}
