package functional

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

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
			"Last Verification",
		)

		for _, project := range projects {
			table.Row(
				project.Id,
				project.ProjectName,
				project.Folder,
				project.LastVerification,
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

	project.LastVerification = time.Now().Format(time.RFC1123)

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

// VerifyProjects verifies all projects in the list of projects.json file.
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

// AddProject adds a project to the list listing all its files.
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
	project.LastVerification = time.Now().Format(time.RFC1123)
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

	// verify existence of c_projectname.cnf and create or not the file
	if _, err := os.Stat(project.Folder + string(os.PathSeparator) + "c_" + project.ProjectName + ".cnf"); err == nil {
		result = append(result, "\n Connection file found!.\n")
	} else {
		err = utils.WriteSampleCNF(project.Folder, project.ProjectName)
		if err != nil {
			return nil, err
		}
		if verbose {
			result = append(result, "\n Connection file sample created.\n")
		}
	}

	// verify existence of projectname.bsh and create or not the file
	if _, err := os.Stat(project.Folder + string(os.PathSeparator) + project.ProjectName + ".bsh"); err == nil {
		result = append(result, "\n External variables file found!.\n")
	} else {
		err = utils.WriteSampleBSH(project.Folder, project.ProjectName)
		if err != nil {
			return nil, err
		}
		if verbose {
			result = append(result, "\n External variables file sample created.\n")
		}
	}

	result = append(result, "\n Project "+name+" saved on projects file!\n")

	return result, nil
}

// UpdateAProject updates a single project. By "update" I mean load all the project files
// to the database selected, but just the marked with "new" or "modified" (but force variable make the flags be ignored).
func UpdateAProject(name string, verbose bool, force bool) ([]string, error) {
	var result []string //TODO: Make verbose

	// Lists the projects
	projects, err := utils.LoadProjectFile()
	if err != nil {
		return nil, err
	}

	// Gets correct project
	idx := slices.IndexFunc(projects, func(c utils.Project) bool { return c.ProjectName == name })

	project := &projects[idx]
	files := &project.Files

	projectLastVerification, err := time.Parse(time.RFC1123, project.LastVerification)
	if err != nil {
		return nil, err
	}

	if projectLastVerification.Sub(time.Now()).Minutes() >= utils.MaxTimeMinutesVerification {
		return nil, errors.New("The last verification of the project is greater than the max time of verification (" +
			strconv.Itoa(utils.MaxTimeMinutesVerification) + " minutes)\n\n Please run beryl vr -n " + project.ProjectName)
	}

	// Gets connection string
	connection, err := utils.ReadCNF(project.Folder + string(os.PathSeparator) + "c_" + project.ProjectName + ".cnf")
	if err != nil {
		return nil, err
	}

	// Gets external variables
	variables, err := utils.ReadExternalVariablesFile(project.Folder + string(os.PathSeparator) + project.ProjectName + ".bsh")
	if err != nil {
		return nil, err
	}

	// Tests the connection
	err = utils.TestConnection(connection)
	if err != nil {
		return nil, err
	}

	// Sends the files
	var counter int = -1
	var removeList []int
	for _, file := range *files {
		counter += 1
		if file.Excluded {
			//TODO: remove excluded files from struct
			removeList = append(removeList, counter)
			continue
		}

		if !force && ((!file.Modified) || (!file.NewFile)) {
			continue
		}

		result, err = utils.SendCodeToDatabase(file.FilePath, variables, connection, verbose)
		if err != nil {
			return nil, err
		}
		file.Modified = false
		file.NewFile = false
	}

	//TODO: save struct after all updates
	for _, item := range removeList {
		utils.RemoveItemFromFiles(*files, item)
	}
	utils.SaveProjectFile(projects)

	result = append(result, "Project "+project.ProjectName+" updated!")

	return result, nil
}

func UpdateProjects(verbose bool, force bool) ([]string, error) {
	projects, err := utils.LoadProjectFile()
	if err != nil {
		return nil, err
	}

	var result []string = nil

	for _, project := range projects {
		text, err := UpdateAProject(project.ProjectName, verbose, force)
		if err != nil {
			return nil, err
		}

		result = append(result, text...)
	}
	return result, nil
}

func TestAConnection(name string) ([]string, error) {
	var result []string

	// Lists the projects
	projects, err := utils.LoadProjectFile()
	if err != nil {
		return nil, err
	}

	// Gets correct project
	idx := slices.IndexFunc(projects, func(c utils.Project) bool { return c.ProjectName == name })

	project := &projects[idx]

	connection, err := utils.ReadCNF(project.Folder)
	if err != nil {
		return nil, err
	}

	err = utils.TestConnection(connection)
	if err != nil {
		return nil, err
	}

	result = append(result, "Connection with "+project.ProjectName+" worked!")

	return result, nil
}
