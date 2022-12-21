package utils

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"
)

// SaveProjectFile saves the Project struct inside a json file.
func SaveProjectFile(data []Project) {
	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Fatalf("could not convert struct data into json file: %v", err)
	}

	//Save Json Data  into a json file
	if err = os.WriteFile("projects.json", jsonData, 0644); err != nil {
		log.Fatalf("could not saveJSON file: %v", err)
	}
}

// LoadProjectFile loads a json file inside the Project struct
func LoadProjectFile() ([]Project, error) {
	file, err := os.Open("projects.json")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	jsonData := make([]byte, 1024)
	_, err = file.Read(jsonData)
	if err != nil {
		return nil, err
	}

	var data []Project

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}

	return data, nil

}

func ListFilesInFolder(startPath string, newFile bool) ([]File, error) {
	var list []File
	// Get a list of all files in the root folder and its subfolders.
	files, err := filepath.Glob(startPath + "/**/*.sql")
	if err != nil {
		return nil, err
	}

	// Iterate through the list of files and print their names.
	for _, file := range files {
		// Get the file info.
		info, err := os.Stat(file)
		if err != nil {
			log.Fatal(err)
		}
		fileItem := File{
			FilePath:         file,
			LastModification: info.ModTime().Format(time.RFC1123),
			Modified:         newFile,
			Exists:           newFile,
			NewFile:          newFile,
			Excluded:         false,
		}

		list = append(list, fileItem)
	}
	return list, nil
}
