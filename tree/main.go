package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var ignoreFolders = []string{"node_modules", "vendor", ".git", "dist", "build"}

// Node represents a file or directory
type Node struct {
	Name     string
	Path     string
	IsDir    bool
	Children []Node
}

// shouldIgnore checks if a folder should be ignored
func shouldIgnore(name string) bool {
	for _, ignored := range ignoreFolders {
		if name == ignored {
			return true
		}
	}
	return false
}

// BuildTree recursively builds the directory tree
func BuildTree(path string) (Node, error) {
	info, err := os.Stat(path)
	if err != nil {
		return Node{}, err
	}

	node := Node{
		Name:  info.Name(),
		Path:  path,
		IsDir: info.IsDir(),
	}

	if !info.IsDir() {
		return node, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return node, nil
	}

	for _, entry := range entries {
		// Skip hidden files and directories
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		// Skip ignored folders
		if entry.IsDir() && shouldIgnore(entry.Name()) {
			continue
		}

		childPath := filepath.Join(path, entry.Name())
		child, err := BuildTree(childPath)
		if err != nil {
			continue
		}

		node.Children = append(node.Children, child)
	}

	// Sort children: directories first, then by name
	sort.Slice(node.Children, func(i, j int) bool {
		if node.Children[i].IsDir != node.Children[j].IsDir {
			return node.Children[i].IsDir
		}
		return node.Children[i].Name < node.Children[j].Name
	})

	return node, nil
}

// PrintTree prints the tree structure with proper formatting
func PrintTree(node Node, prefix string, isLast bool, isRoot bool) {
	// Print current node
	if isRoot {
		// Root node
		fmt.Println(node.Name + "/")
		fmt.Println("│")
	} else {
		// Branch formatting
		connector := "├── "
		if isLast {
			connector = "└── "
		}

		if node.IsDir {
			fmt.Printf("%s%s%s/\n", prefix, connector, node.Name)
		} else {
			fmt.Printf("%s%s%s\n", prefix, connector, node.Name)
		}
	}

	// Print children
	if node.IsDir && len(node.Children) > 0 {
		for i, child := range node.Children {
			isLastChild := i == len(node.Children)-1

			var newPrefix string
			if isRoot {
				newPrefix = ""
			} else {
				if isLast {
					newPrefix = prefix + " "
				} else {
					newPrefix = prefix + "│"
				}
			}

			PrintTree(child, newPrefix, isLastChild, false)

			// Add empty line after each child (except the last one)
			if !isLastChild {
				if isRoot {
					fmt.Println("│")
				} else {
					fmt.Printf("%s\n", newPrefix)
				}
			}
		}
	}
}

func main() {
	// Define command-line flag for ignore folders
	ignoreList := flag.String("ignore", "", "Comma-separated list of folders to ignore")
	flag.Parse()

	// Update ignore list if provided
	if *ignoreList != "" {
		customIgnores := strings.Split(*ignoreList, ",")
		for i := range customIgnores {
			customIgnores[i] = strings.TrimSpace(customIgnores[i])
		}
		ignoreFolders = append(ignoreFolders, customIgnores...)
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Build the tree
	root, err := BuildTree(cwd)
	if err != nil {
		fmt.Printf("Error building tree: %v\n", err)
		os.Exit(1)
	}

	// Print the tree
	PrintTree(root, "", true, true)
}
