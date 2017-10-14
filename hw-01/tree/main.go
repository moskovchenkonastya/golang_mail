package main


import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

func readDir(path string, printFiles bool) ([]string, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdir(-1)
	f.Close();
	if err != nil {
		return nil, err
	}
	
	values := []string{}
	
	for _, file := range names {

		if file.Name() == ".DS_Store" {
			continue
		}
		if printFiles || file.IsDir() {
			values = append(values, file.Name())
		}
	}
	sort.Strings(values) 
	return values, err
}

func getSizeDesc(fullFile string) string {
	stat,_ := os.Stat(fullFile)
	size := stat.Size()
	if !stat.IsDir() {
		if fullFile == "main.go" {
			return " (vary)"
		} else if size > 0 {
			return " (" + strconv.FormatInt(size, 10) + "b)"
		} else {
			return " (empty)"
		}
	}	
	return ""
}

//рекурсивный обход папок дерева
func dirTree(out io.Writer, path string, printFiles bool) error{
	sep := ""
	return dirTreeSep(out, path, printFiles, sep)
}

// передача отсутстпов в виде строки
func dirTreeSep(out io.Writer, path string, printFiles bool, sep string) error{
	// формирование массива близлежащих элементов в одной дирректории
	names, _ := readDir(path, printFiles)
	dirCount := len(names)
	var  curSep, nextSep string
	
	for i, name := range names {
		cur := filepath.Join(path, name)

		if i < dirCount - 1 {
			curSep = "├───"	
		} else {
			curSep = "└───"	
		}
		if i < dirCount - 1 {
			nextSep = "│	"	
		} else {
			nextSep = "\t"	
		}
		if printFiles{
			fmt.Fprint(out, sep + curSep, name, getSizeDesc(cur))
			fmt.Fprint(out,"\n")
		} else {
			fmt.Fprint(out, sep + curSep, name)
			fmt.Fprint(out,"\n")
		}
		dirTreeSep(out, cur, printFiles, sep + nextSep)
	}

	return nil
}

func main() {
	
	  out := os.Stdout
	  if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	  }
	  path := os.Args[1]
	  printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	  dirTree(out, path, printFiles) 
}
	