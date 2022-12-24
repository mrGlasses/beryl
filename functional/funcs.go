package functional

import (
	"errors"
	"fmt"
	"sort"

	"github.com/ecoshub/stable"
	"github.com/mrGlasses/beryl/utils"
	"golang.org/x/exp/slices"
)

// ListProjectData shows a table with the projects (name == "") or the files of the project (name == any valid name).
func ListProjectData(name string) (string, error) {
	projects, err := utils.LoadProjectFile()
	if err != nil {
		return "", err
	}

	// var table *stable.STable
	if name == "" {
		table := stable.New("List of projects")
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
		return table.String(), nil
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
			"New File",
			"Excluded",
		)

		files := projects[idx].Files

		for _, line := range files {
			table.Row(
				line.FilePath,
				line.LastModification,
				line.Modified,
				line.NewFile,
				line.Excluded,
			)
		}
		return table.String(), nil
	}
}

// VerifyAProject verifies a project files by his given name whenether
// a file is modified, excluded or added.
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

	verifier, _, _ := utils.ListFilesInFolder(project.Folder, true, false)

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

	//granting the files order
	projectReorder := projects[idx]
	filesReorder := projectReorder.Files
	sort.Slice(filesReorder, func(i, j int) bool {
		return filesReorder[i].FilePath < filesReorder[j].FilePath
	},
	)

	files = &filesReorder

	utils.SaveProjectFile(projects)

	result = append(result, "")
	final := fmt.Sprintf("Final Files Status: Modified: %v - Added: %v - Deleted: %v", modifications.Modified, modifications.New, modifications.Excluded)
	result = append(result, final)
	return result, nil
}

// VerifyProjects verify all projects in the list of projects.json file.
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

func AddProject(name string, location string, verbose bool) ([]string, error) {
	var result []string
	var speak []string
	var err error

	projects, err := utils.LoadProjectFile()
	if err != nil {
		//it's the first one
		if verbose {
			result = append(result, "\n First project!\n")
		}
		err = nil
	}
	var project utils.Project

	//project creation and file loading
	project.Id = utils.GetLastID(projects)
	project.ProjectName = name
	project.Folder = location
	project.Files, speak, err = utils.ListFilesInFolder(location, true, verbose)
	if err != nil {
		return nil, err
	}

	result = append(result, speak...)

	if verbose {
		result = append(result, "\n Files loaded.\n")
	}
	//add project to the slice
	projects = append(projects, project)

	//save project
	err = utils.SaveProjectFile(projects)
	if err != nil {
		return nil, err
	}
	if verbose {
		result = append(result, "\n Project "+name+" saved on projects file!\n")
	}

	return result, nil
}
