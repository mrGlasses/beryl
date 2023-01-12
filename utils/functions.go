package utils

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "gopkg.in/rana/ora.v3"
)

// SaveProjectFile saves the Project struct inside a json file.
func SaveProjectFile(data []Project) error {
	var result error = nil

	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		result = fmt.Errorf("could not convert struct data into json file: %v", err)
	}

	//Save Json Data  into a json file
	homedir, _ := os.UserHomeDir()

	if err = os.WriteFile(homedir+string(os.PathSeparator)+"projects.json", jsonData, 0644); err != nil {
		result = fmt.Errorf("could not save JSON file: %v", err)
	}

	return result
}

// LoadProjectFile loads a json file inside the Project struct.
func LoadProjectFile() ([]Project, error) {
	homedir, _ := os.UserHomeDir()
	_, err := os.Stat(homedir + string(os.PathSeparator) + "projects.json")
	if err != nil {
		return nil, errors.New("no projects added")
	}

	file, err := os.ReadFile(homedir + string(os.PathSeparator) + "projects.json")
	if err != nil {
		return nil, err
	}

	var data []Project

	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// ListFilesInFolder list all *.sql files in a given folder and subfolders
// returning a slice of File(struct) and marking them as new or not using
// the second parameter.
func ListFilesInFolder(startPath string, newPath bool, verbose bool) ([]*File, []string, error) {
	var list []*File
	var speak []string

	// Get a list of all files in the root folder and its subfolders.
	// folderSearch := startPath + string(os.PathSeparator) + "**" + string(os.PathSeparator) + "*.sql"
	
	// files, err := filepath.Glob(folderSearch)
	// if err != nil {
	// 	return nil, nil, err
	// }

	// Get a list of all files in the root folder and its subfolders.
	err := filepath.Walk(startPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() && filepath.Ext(path) == ".sql" {
			fileItem := File{
				FilePath:         path,
				LastModification: info.ModTime().Format(time.RFC1123),
				Modified:         newPath,
				NewFile:          newPath,
				Excluded:         false,
				// Exists:           newPath,
			}
			if verbose {
				speak = append(speak, "File "+path+" added!")
			}
			list = append(list, &fileItem)
        }
        return nil
    })

    if err != nil {
        return nil, nil, err
    }


	// Iterate through the list of files and print their names.
	// for _, file := range files {
	// 	// Get the file info.
	// 	info, err := os.Stat(file)
	// 	if err != nil {
	// 		return nil, nil, err
	// 	}
	// 	fileItem := File{
	// 		FilePath:         file,
	// 		LastModification: info.ModTime().Format(time.RFC1123),
	// 		Modified:         newPath,
	// 		NewFile:          newPath,
	// 		Excluded:         false,
	// 		// Exists:           newPath,
	// 	}
	// 	if verbose {
	// 		speak = append(speak, "File "+file+" added!")
	// 	}
	// 	list = append(list, &fileItem)
	// }
	return list, speak, nil
}

// GetLastID gets the last ID from a project in a given slice of Project(struct).
func GetLastID(projects []Project) int {

	if len(projects) == 0 {
		return 1
	}

	last := len(projects) - 1

	return projects[last].Id
}

// ReadCNF reads the CNF file (c_<ProjectName>.cnf) in the main folder of the project
// and then return a connection string struct.
func ReadCNF(filePath string) (ConnString, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		// Handle the error
		empty := new(ConnString)
		return *empty, err
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Initialize variables to store the connection settings
	var host string
	var driver string
	var port string
	var user string
	var password string
	var database string

	// Read the file line by line
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line is a comment
		if strings.HasPrefix(line, "c") {
			continue
		}

		// Split the line on the "=" character to get the key and value
		parts := strings.Split(line, "=")
		key := parts[0]
		value := parts[1]

		// Set the appropriate variable based on the key
		switch key {
		case "driver":
			driver = value
		case "host":
			host = value
		case "port":
			port = value
		case "user":
			user = value
		case "password":
			password = value
		case "database":
			database = value
		}
	}

	// Check for any scanning errors
	if err := scanner.Err(); err != nil {
		// Handle the error
		empty := new(ConnString)
		return *empty, err
	}

	// Return the connection settings
	result := ConnString{
		DbsName:  driver,
		User:     user,
		Password: password,
		Server:   host,
		Database: database,
		Port:     port,
	}

	return result, nil
}

// ReadExternalVariablesFile reads the variables file (<ProjectName>.bsh) in the main folder of the project
// and then return a slice of external variables struct.
func ReadExternalVariablesFile(filePath string) ([]ExternalVariables, error) {
	var result []ExternalVariables

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		// Handle the error
		return nil, err
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Read the file line by line
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line is a comment
		if strings.HasPrefix(line, "c") {
			continue
		}

		// Split the line on the ";" character to get the key and value
		parts := strings.Split(line, ";")
		code := parts[0]
		rplc := parts[1]

		pair := ExternalVariables{
			Old: []byte(code),
			New: []byte(rplc),
		}

		result = append(result, pair)

	}

	// Check for any scanning errors
	if err := scanner.Err(); err != nil {
		// Handle the error
		return nil, err
	}

	// Return the list

	return result, nil
}

// SendCodeToDatabase gets the *.sql file, replace the variables in the slice of struct ([]External Variables)
// and send to the database using the ConnectionString struct.
func SendCodeToDatabase(filePath string, extVar []ExternalVariables, connection ConnString, verbose bool) ([]string, error) {
	var result []string
	var conString string

	// Select the connection
	switch connection.DbsName {
	case "mysql":
		conString = connection.User + ":" + connection.Password + "@tcp(" + connection.Server + ":" + connection.Port + ")/" + connection.Database
	case "ora":
		conString = connection.User + "/" + connection.Password + "@" + connection.Server + ":" + connection.Port + "/" + connection.Database
	case "mssql":
		conString = fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s", connection.Server+":"+connection.Port, connection.User, connection.Password, connection.Database)
	case "postgres":
		conString = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", connection.User, connection.Password, connection.Server, connection.Port, connection.Database)
	}

	db, err := sql.Open(connection.DbsName, conString)
	if err != nil {
		return nil, fmt.Errorf(filePath+ ": "+ err.Error())
	}
	defer db.Close()

	// Read the SQL file into a slice of bytes
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf(filePath+ ": "+ err.Error())
	}

	if verbose {
		result = append(result, "File read.")
	}

	// Replacing variables
	for _, pair := range extVar {
		data = bytes.ReplaceAll(data, []byte(pair.Old), []byte(pair.New))
	}

	if verbose {
		result = append(result, "Variables replaced.")
	}

	// Execute the commands in the SQL file
	_, err = db.Exec(string(data))
	if err != nil {
		return nil, fmt.Errorf(filePath+ ": "+ err.Error())
	}
	result = append(result, "File "+filePath+" was executed successfully!")
	return result, nil
}

// TestConnection tests a connection of a given connection string.
func TestConnection(connection ConnString) error {
	var conString string

	switch connection.DbsName {
	case "mysql":
		conString = connection.User + ":" + connection.Password + "@tcp(" + connection.Server + ":" + connection.Port + ")/" + connection.Database
	case "ora": //BECAUSE OF THIS SHITHEAD YOU MUST INSTALL GCC AND OCI8
		conString = connection.User + "/" + connection.Password + "@" + connection.Server + ":" + connection.Port + "/" + connection.Database
	case "mssql":
		conString = fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s", connection.Server+":"+connection.Port, connection.User, connection.Password, connection.Database)
	case "postgres":
		conString = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", connection.User, connection.Password, connection.Server, connection.Port, connection.Database)
	}

	db, err := sql.Open(connection.DbsName, conString)
	if err != nil {
		return err
	}

	defer db.Close()
	return err
}

// WriteSampleCNF creates the basic conection file for a project in a project folder.
func WriteSampleCNF(filePath string, projectName string) error {
	var lines []string

	lines = append(lines, CnfSampleLn1)
	lines = append(lines, CnfSampleLn2)
	lines = append(lines, CnfSampleLn3)
	lines = append(lines, CnfSampleLn4)
	lines = append(lines, CnfSampleLn5)
	lines = append(lines, CnfSampleLn6)

	// Try to open/create file
	file, err := os.OpenFile(filePath+string(os.PathSeparator)+"c_"+projectName+".cnf", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a buffered writer
	writer := bufio.NewWriter(file)

	// Write the lines to the file
	for _, line := range lines {
		_, err = writer.WriteString(line)
		if err != nil {
			return err
		}
	}

	// Flush the buffer to the file
	writer.Flush()

	return nil
}

// WriteSampleBSH creates the basic external variables file for a project in a project folder.
func WriteSampleBSH(filePath string, projectName string) error {
	var lines []string

	lines = append(lines, BshSampleLn1)
	lines = append(lines, BshSampleLn2)

	// Try to open/create file
	file, err := os.OpenFile(filePath+string(os.PathSeparator)+projectName+".bsh", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a buffered writer
	writer := bufio.NewWriter(file)

	// Write the lines to the file
	for _, line := range lines {
		_, err = writer.WriteString(line)
		if err != nil {
			return err
		}
	}

	// Flush the buffer to the file
	writer.Flush()

	return nil
}

// removeItemFromFiles remove a item from the []File slice.
func RemoveItemFromFiles(slice []*File, idx int) []*File {
	copy(slice[idx:], slice[idx+1:])
	empty := new(File)
	slice[len(slice)-1] = empty
	slice = slice[:len(slice)-1]
	return slice
}

// RemoveItemFromProjects remove a item from the []Project slice.
func RemoveItemFromProjects(slice []Project, idx int) []Project {
	copy(slice[idx:], slice[idx+1:])
	empty := new(Project)
	slice[len(slice)-1] = *empty
	slice = slice[:len(slice)-1]
	return slice
}

// RenameProject just recieve a project and change the ProjectName property for a new name.
func RenameProject(project *Project, name string) {
	project.ProjectName = name
}

// ReplaceProject just recieve a project and change the Folder property for a new folder.
func ReplaceProject(project *Project, folder string) {
	project.Folder = folder
}

// CopyFile - https://github.com/mactsouk/opensource.com/blob/master/cp3.go
func CopyFile(src, dst string, BUFFERSIZE int64) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	_, err = os.Stat(dst)
	if err == nil {
		return fmt.Errorf("file %s already exists", dst)
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	if err != nil {
		panic(err)
	}

	buf := make([]byte, BUFFERSIZE)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}

// Confirm - https://gist.github.com/r0l1/3dcbb0c8f6cfe9c66ab8008f55f8f28b?permalink_comment_id=1995588#gistcomment-1995588
func Confirm(s string) bool {
	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		res, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		// Empty input (i.e. "\n")
		if len(res) < 2 {
			continue
		}

		// Wrong answer
		result := strings.ToLower(strings.TrimSpace(res))[0]
		if (result != 'y') && (result != 'n') {
			continue
		}

		return result == 'y'
	}
}
