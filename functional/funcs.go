package functional

import (
	"errors"
	"fmt"

	"github.com/ecoshub/stable"
	"github.com/mrGlasses/BerylSQLHelper/utils"
	"golang.org/x/exp/slices"
)

// ListProjectData shows a table with the projects (name == "") or the files of the project (name == any valid name)
func ListProjectData(name string) (string, error) {
	projects, err := utils.LoadProjectFile()
	if err != nil {
		return "", err
	}

	var table stable.STable
	if name == "" {
		table := stable.New("List of Projects")
		table.AddFields(
			"ID", //TODO: Ajustar posição do ID (add with options)
			"Name",
			"Folder",
		)

		for _, project := range projects {
			table.Row(
				project.Id,
				project.ProjectName,
				project.Folder,
			)
		}
	} else {
		idx := slices.IndexFunc(projects, func(c utils.Project) bool { return c.ProjectName == name })
		if idx == -1 {
			return "", errors.New("project not found")
		}

		table := stable.New("Project: " + projects[idx].ProjectName)
		table.AddFields(
			"File Path",
			"Last Modification",
			"Modified",
			// "Exists",
			"New File",
			"Excluded",
		)

		files := projects[idx].Files

		for _, line := range files {
			table.Row(
				line.FilePath,
				line.LastModification,
				line.Modified,
				// line.Exists,
				line.NewFile,
				line.Excluded,
			)
		}
	}
	return table.String(), nil
}

// VerifyAProject
func VerifyAProject(name string, verbose bool) ([]string, error) {
	var modifications utils.FileStatus
	modifications.ProjectName = name
	modifications.New = 0
	modifications.Modified = 0
	modifications.Excluded = 0

	var result []string

	projects, err := utils.LoadProjectFile()
	if err != nil {
		return nil, err
	}

	idx := slices.IndexFunc(projects, func(c utils.Project) bool { return c.ProjectName == name })

	project := &projects[idx]
	files := &project.Files

	verifier, _ := utils.ListFilesInFolder(project.Folder, true)

	for _, file := range *files {
		for _, fileV := range verifier {
			// Exists
			if file.FilePath == fileV.FilePath {
				// Modified
				if file.LastModification != fileV.LastModification {
					if verbose {
						result = append(result, "MODIFIED: "+file.FilePath)
					}
					file.LastModification = fileV.LastModification
					file.Modified = true
					modifications.Modified += 1
				}
				continue
			}
		}
		if verbose {
			result = append(result, "EXCLUDED: "+file.FilePath)
		}
		file.Excluded = true
		modifications.Excluded += 1
	}

	// For new entries
	for _, fileV := range verifier {
		exists := false

		for _, file := range *files {
			if file.FilePath == fileV.FilePath {
				exists = true
				break
			}
		}
		if exists {
			continue
		}
		*files = append(*files, fileV)
		modifications.New += 1
		if verbose {
			result = append(result, "ADDED: "+fileV.FilePath)
		}
	}

	utils.SaveProjectFile(projects)

	result = append(result, "")
	final := fmt.Sprintf("Final Files Status: Modified: %v - Added: %v - Deleted: %v", modifications.Modified, modifications.New, modifications.Excluded)
	result = append(result, final)
	return result, nil
}

func VerifyProjects(verbose bool) ([]string, error) {
	projects, err := utils.LoadProjectFile()
	if err != nil {
		return nil, err
	}

	var result []string = nil

	for _, project := range projects {
		text, err := VerifyAProject(project.ProjectName, verbose)
		if err != nil {
			return nil, err
		}

		result = append(result, text...)
	}
	return result, nil
}
