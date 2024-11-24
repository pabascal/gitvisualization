package main

import (
		"bufio"
		"fmt"
		"io"
		"log"
		"os"
		"os/user"
		"strings"
)


func getDotFilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	dotFile := usr.HomeDir + "/.gogitlocalstats"  
	//usr = current user, and usr.HomeDir is the "path to the user's home directory"

	return dotFile
}

func openFile(filePath string) *os.File {
	// 0755 argument sets file permissions (owner can read,write, execute, others can read, execute)
	// O_Append - open file in append mode. O_Wronly - open file for writing only
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_RDWR, 0755)  //or O_WRONLY
	if err != nil {
		if os.IsNotExist(err) {
			_, err = os.Create(filePath)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	return f
}

func parseFileLinesToSlice(filePath string) []string {
	f := openFile(filePath)
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			panic(err)

		}
	}
	return lines
}

func sliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

//joinSlices loops through new repos and appends them all to the existing repos slice (joining)
func joinSlices(new []string, existing []string) []string {
	for _, i := range new {
		if !sliceContains(existing, i) {
			existing = append(existing, i)
		}
	}
	return existing
}

func dumpStringsSliceToFile(repos []string, filePath string) {
	content := strings.Join(repos, "\n")
	os.WriteFile(filePath, []byte(content), 0755)
}

// addNewSliceElementsToFile given a slice of strings representing paths, stores them
// to the filesystem
func addNewSliceElementsToFile(filePath string, newRepos []string) {
	existingRepos := parseFileLinesToSlice(filePath)
	repos := joinSlices(newRepos, existingRepos)
	dumpStringsSliceToFile(repos, filePath)

}


func recursiveScanFolder(folder string) []string {
	return scanGitFolders(make([]string, 0), folder)

}

func scan(folder string) {
	fmt.Printf("Found folders:\n\n")
	repositories := recursiveScanFolder(folder)
	filePath := getDotFilePath()
	addNewSliceElementsToFile(filePath, repositories)
	fmt.Printf("\n\nSuccessfully added\n\n")
}

//Recursively searches through all folders within the input directory location
//Reutrns a list of all .git folders, returns the base folder/parent folder of the .git folder
func scanGitFolders(folders []string, folder string) []string {
	
	//Trims off the ending / in the file path if it exists
	folder = strings.TrimSuffix(folder, "/")

	f, err := os.Open(folder)
	if err != nil {
		log.Fatal(err)       //err contains potential error from os.Open call. 
		                     //if err exists (is not nil), log.Fatal(err) will terminate
							 //program and print error message contained in err
	}

	files, err := f.Readdir(-1)     //reads contents of folder f, stores the list of contents in
	                                //'files'. -1 means it will read all files.
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	var path string

	for _, file := range files {
		if file.IsDir() {         // IsDir evaluates if it's a folder/directory that can be further opened
			path = folder + "/" + file.Name()
			if file.Name() == ".git" {
				path = strings.TrimSuffix(path, "/.git")
				fmt.Println(path)     //print location path where a .git is contained
				folders = append(folders, path)
				continue
			}
			if file.Name() == "vendor" || file.Name() == "node_modules" {
				continue
			}
			folders = scanGitFolders(folders, path)
		}
	}
	return folders
}


