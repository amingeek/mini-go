package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	books := make(map[string]string)
	var isbns []string

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())

	for i := 0; i < n; i++ {
		scanner.Scan()
		parts := strings.Fields(scanner.Text())

		if parts[0] == "ADD" {
			books[parts[1]] = strings.Join(parts[2:], " ")
		} else if parts[0] == "REMOVE" {
			delete(books, parts[1])
		}
	}

	for isbn := range books {
		isbns = append(isbns, isbn)
	}

	sort.Slice(isbns, func(i, j int) bool {
		titleI := books[isbns[i]]
		titleJ := books[isbns[j]]
		if titleI == titleJ {
			numI, _ := strconv.Atoi(isbns[i])
			numJ, _ := strconv.Atoi(isbns[j])
			return numI < numJ
		}
		return titleI < titleJ
	})

	for _, isbn := range isbns {
		fmt.Println(isbn)
	}
}
