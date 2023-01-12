package functional

import (
	"errors"
	"fmt"
	"math"
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
			return "", fmt.Errorf("project %s not found", name)
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
	if idx == -1 {
		return nil, fmt.Errorf("project %s not found", name)
	}

	project := &projects[idx]
	files := &project.Files

	verifier, _, _ := utils.ListFilesInFolder(project.Folder, true, verbose)
	exists := false

	for _, file := range project.Files {
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
				exists = true
				break
			}
		}
		if exists {
			exists = false
			continue
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
		result = append(result, "\n##### PROJECT: "+project.ProjectName)
		text, err := VerifyAProject(project.ProjectName, verbose)
		if err != nil {
			return nil, err
		}

		result = append(result, text...)
		result = append(result, "____________________________________________\n")
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

	exists := slices.IndexFunc(projects, func(c utils.Project) bool { return c.ProjectName == name })
	if exists != -1 {
		return nil, fmt.Errorf("the project with name %s already exists", name)
	}

	var project utils.Project

	//project creation and file loading
	project.Id = utils.GetLastID(projects) + 1
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
	var result []string

	// Lists the projects
	projects, err := utils.LoadProjectFile()
	if err != nil {
		return nil, err
	}

	// Gets correct project
	idx := slices.IndexFunc(projects, func(c utils.Project) bool { return c.ProjectName == name })
	if idx == -1 {
		return nil, fmt.Errorf("project %s not found", name)
	}

	project := &projects[idx]

	files := &project.Files
	if verbose {
		result = append(result, "Project found.")
	}

	projectLastVerification, err := time.Parse(time.RFC1123, project.LastVerification)
	if err != nil {
		return nil, err
	}

	if math.Abs(time.Until(projectLastVerification).Minutes()) >= utils.MaxTimeMinutesVerification {
		return nil, errors.New("the last verification of the project is greater than the max time of verification (" +
			strconv.Itoa(utils.MaxTimeMinutesVerification) + " minutes)\n\n Please run beryl vr -n " + project.ProjectName)
	}

	if verbose {
		result = append(result, "Time of last verification OK.")
	}

	// Gets connection string
	connection, err := utils.ReadCNF(project.Folder + string(os.PathSeparator) + "c_" + project.ProjectName + ".cnf")
	if err != nil {
		return nil, err
	}

	if verbose {
		result = append(result, "Connection string found.")
	}

	// Gets external variables
	variables, err := utils.ReadExternalVariablesFile(project.Folder + string(os.PathSeparator) + project.ProjectName + ".bsh")
	if err != nil {
		return nil, err
	}

	if verbose {
		result = append(result, "Variables found.")
	}

	// Tests the connection
	err = utils.TestConnection(connection)
	if err != nil {
		return nil, err
	}

	if verbose {
		result = append(result, "Connection OK.")
	}

	// Sends the files
	var counter int = -1
	var removeList []int
	for _, file := range *files {
		counter += 1
		if file.Excluded {
			//TODO: remove excluded files from struct
			removeList = append(removeList, counter)

			if verbose {
				result = append(result, "File "+file.FilePath+" to remove list.")
			}
			continue
		}

		if !force && ((!file.Modified) || (!file.NewFile)) {

			if verbose {
				result = append(result, "File "+file.FilePath+" is not a new/modified file.")
			}
			continue
		}

		if verbose {
			result = append(result, "Preparing file "+file.FilePath+" for send to the database...")
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
		*files = utils.RemoveItemFromFiles(*files, item)

		if verbose {
			result = append(result, "File in index "+strconv.Itoa(item)+" removed from project.")
		}
	}

	utils.SaveProjectFile(projects)

	if verbose {
		result = append(result, "Project saved.")
	}

	result = append(result, "Project "+project.ProjectName+" updated!")

	return result, nil
}

// UpdateProjects gets the list of all projects then update each one.
func UpdateProjects(verbose bool, force bool) ([]string, error) {
	projects, err := utils.LoadProjectFile()
	if err != nil {
		return nil, err
	}

	var result []string = nil

	for _, project := range projects {
		result = append(result, "\n##### PROJECT: "+project.ProjectName)
		text, err := UpdateAProject(project.ProjectName, verbose, force)
		if err != nil {
			return nil, err
		}

		result = append(result, text...)
		result = append(result, "____________________________________________\n")
	}
	return result, nil
}

// Test the connection using the connection file (c_projectName.cnf) with the database.
func TestAConnection(name string) ([]string, error) {
	var result []string

	// Lists the projects
	projects, err := utils.LoadProjectFile()
	if err != nil {
		return nil, err
	}

	// Gets correct project
	idx := slices.IndexFunc(projects, func(c utils.Project) bool { return c.ProjectName == name })
	if idx == -1 {
		return nil, fmt.Errorf("project %s not found", name)
	}

	project := &projects[idx]

	connection, err := utils.ReadCNF(project.Folder + string(os.PathSeparator) + "c_" + project.ProjectName + ".cnf")
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

// Rename a project using the ID
func RenameAProject(id int, newName string) ([]string, error) {
	var result []string

	// Lists the projects
	projects, err := utils.LoadProjectFile()
	if err != nil {
		return nil, err
	}

	// Gets correct project
	idx := slices.IndexFunc(projects, func(c utils.Project) bool { return c.Id == id })

	// Project found
	if idx != -1 {
		if !utils.Confirm(fmt.Sprintf("Are you sure you want to rename the project to %s?\n", newName)) {
			return nil, errors.New("rename cancelled by user")
		}
	} else {
		return nil, fmt.Errorf("project id %d not found", id)
	}

	project := &projects[idx]

	exists := slices.IndexFunc(projects, func(c utils.Project) bool { return c.ProjectName == newName })
	if exists != -1 {
		return nil, fmt.Errorf("the project with name %s already exists", newName)
	}

	oldName := project.ProjectName

	utils.RenameProject(project, newName)

	// Renaming files
	err = os.Rename(project.Folder+string(os.PathSeparator)+"c_"+oldName+".cnf", project.Folder+string(os.PathSeparator)+"c_"+project.ProjectName+".cnf")
	if err != nil {
		return nil, err
	}

	err = os.Rename(project.Folder+string(os.PathSeparator)+oldName+".bsh", project.Folder+string(os.PathSeparator)+project.ProjectName+".bsh")
	if err != nil {
		return nil, err
	}

	err = utils.SaveProjectFile(projects)
	if err != nil {
		return nil, err
	}

	result = append(result, "Project renamed!")

	return result, nil
}

// ReplaceAProject just replaces the main folder in the project then re-verifies everything.
func ReplaceAProject(name string, newFolder string, verbose bool) ([]string, error) {
	var result []string

	// Lists the projects
	projects, err := utils.LoadProjectFile()
	if err != nil {
		return nil, err
	}

	// Gets correct project
	idx := slices.IndexFunc(projects, func(c utils.Project) bool { return c.ProjectName == name })

	// Project found
	if idx != -1 {
		if !utils.Confirm(fmt.Sprintf("Are you sure you want to change the folder of the project to %s?\n", newFolder)) {
			return nil, errors.New("folder change cancelled by user")
		}
	} else {
		return nil, fmt.Errorf("project %s not found", name)
	}
	project := &projects[idx]

	// oldFolder := project.Folder

	utils.ReplaceProject(project, newFolder)

	// Copy files - may be useless, because need to remap all files, but I'll let this here.

	// err = utils.CopyFile(oldFolder+string(os.PathSeparator)+"c_"+project.ProjectName+".cnf",
	// 	project.Folder+string(os.PathSeparator)+"c_"+project.ProjectName+".cnf",
	// 	1000000)
	// if err != nil {
	// 	if err.Error() == fmt.Sprintf("File %s already exists.", project.Folder+string(os.PathSeparator)+"c_"+project.ProjectName+".cnf") {
	// 		result = append(result, err.Error())
	// 	} else {
	// 		return nil, err
	// 	}
	// }

	// err = utils.CopyFile(oldFolder+string(os.PathSeparator)+project.ProjectName+".bsh",
	// 	project.Folder+string(os.PathSeparator)+project.ProjectName+".bsh",
	// 	1000000)
	// if err != nil {
	// 	if err.Error() == fmt.Sprintf("File %s already exists.", project.Folder+string(os.PathSeparator)+project.ProjectName+".bsh") {
	// 		result = append(result, err.Error())
	// 	} else {
	// 		return nil, err
	// 	}
	// }

	// Re-verify all project
	vResult, err := VerifyAProject(project.ProjectName, verbose)
	if err != nil {
		return nil, err
	}
	result = append(result, vResult...)

	err = utils.SaveProjectFile(projects)
	if err != nil {
		return nil, err
	}

	result = append(result, "Project replaced!")

	return result, nil
}

// DeleteAProject deletes a project from the projects file.
func DeleteAProject(name string) ([]string, error) {
	var result []string

	// Lists the projects
	projects, err := utils.LoadProjectFile()
	if err != nil {
		return nil, err
	}

	// Gets correct project
	idx := slices.IndexFunc(projects, func(c utils.Project) bool { return c.ProjectName == name })

	// Project found
	if idx != -1 {
		if !utils.Confirm(fmt.Sprintf("Are you sure you want to delete the project called %s?\n", name)) {
			return nil, errors.New("folder change cancelled by user")
		}
	} else {
		return nil, fmt.Errorf("project %s not found", name)
	}

	projects = utils.RemoveItemFromProjects(projects, idx)

	utils.SaveProjectFile(projects)
	result = append(result, "Project deleted!")

	return result, nil
}
