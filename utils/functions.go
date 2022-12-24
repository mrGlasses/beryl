package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
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
	file, err := os.ReadFile(homedir + string(os.PathSeparator) + "projects.json")
	if err != nil {
		return nil, err
	}

	// jsonData := make([]byte, 1024)
	// _, err = file.Read(jsonData)
	// if err != nil {
	// 	return nil, err
	// }

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
func ListFilesInFolder(startPath string, newPath bool, verbose bool) ([]File, []string, error) {
	var list []File
	var speak []string

	// Get a list of all files in the root folder and its subfolders.
	files, err := filepath.Glob(startPath + "/**/*.sql")
	if err != nil {
		return nil, nil, err
	}

	// Iterate through the list of files and print their names.
	for _, file := range files {
		// Get the file info.
		info, err := os.Stat(file)
		if err != nil {
			return nil, nil, err
		}
		fileItem := File{
			FilePath:         file,
			LastModification: info.ModTime().Format(time.RFC1123),
			Modified:         newPath,
			NewFile:          newPath,
			Excluded:         false,
			// Exists:           newPath,
		}

		speak = append(speak, "File "+file+" added!")
		list = append(list, fileItem)
	}
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
