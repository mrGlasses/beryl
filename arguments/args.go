package arguments

import (
	"fmt"
	"os"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/mrGlasses/beryl/functional"
	"github.com/mrGlasses/beryl/utils"
)

func ExecuteArguments(args []string) (string, error) {

	// args = append(args, "u", "-n", "ENNBA","-e") //debug arguments
	fmt.Println(args)

	parser := argparse.NewParser(utils.CommandName, utils.ProgramDescription)

	cmdVersion := parser.Flag("v", "version", &argparse.Options{Required: false, Help: "Shows the installed version of the code"})

	cmdShowAll := parser.NewCommand("sa", "Shows all main folders for each project")

	cmdShow := parser.NewCommand("s", "(Use: beryl s -n projectName) Shows the data of the selected project")
	getShow := cmdShow.String("n", "name", &argparse.Options{Required: true})

	cmdVerifyAll := parser.NewCommand("va", "Verifies all projects and covered folders for updates - -e|--verbose as optional")
	getVAVerbose := cmdVerifyAll.Flag("e", "verbose", &argparse.Options{Required: false})

	cmdVerify := parser.NewCommand("vr", "(Use: beryl vr -n projectName) Verifies a specific project and covered folders for updates - -e|--verbose as optional")
	getVerify := cmdVerify.String("n", "name", &argparse.Options{Required: true})
	getVRVerbose := cmdVerify.Flag("e", "verbose", &argparse.Options{Required: false})

	cmdAddNew := parser.NewCommand("an", "(Use: beryl an -n projectName -l projectLocation) Adds a new project and its folder to the app - -e|--verbose as optional")
	getANName := cmdAddNew.String("n", "name", &argparse.Options{Required: true})
	getANLocation := cmdAddNew.String("l", "location", &argparse.Options{Required: true})
	getANVerbose := cmdAddNew.Flag("e", "verbose", &argparse.Options{Required: false})

	cmdAddHere := parser.NewCommand("ah", "(Use: beryl ah -n projectName) Adds the current folder to the app - -e|--verbose as optional")
	getAHName := cmdAddHere.String("n", "name", &argparse.Options{Required: true})
	getAHVerbose := cmdAddHere.Flag("e", "verbose", &argparse.Options{Required: false})

	cmdUpAll := parser.NewCommand("ua", "(Use: beryl ua )Updates all projects added to the app - -e|--verbose as optional")
	getUAVerbose := cmdUpAll.Flag("e", "verbose", &argparse.Options{Required: false})
	getUAForce := cmdUpAll.Flag("f", "force", &argparse.Options{Required: false, Help: "(Use: [-u projectName|-ua] -f) Only works with -u and -ua command - (be careful) Re-run all files in all folders."})

	cmdUpdate := parser.NewCommand("u", "(Use: beryl u -n projectName) Updates a specific project - -e|--verbose as optional")
	getUpdate := cmdUpdate.String("n", "name", &argparse.Options{Required: true})
	getUVerbose := cmdUpdate.Flag("e", "verbose", &argparse.Options{Required: false})
	getUForce := cmdUpdate.Flag("f", "force", &argparse.Options{Required: false, Help: "(Use: [-u projectName|-ua] -f) Only works with -u and -ua command - (be careful) Re-run all files in all folders."})

	cmdTest := parser.NewCommand("tc", "(Use: beryl tc -n projectName) Test the connection with the server/database")
	getTest := cmdTest.String("n", "name", &argparse.Options{Required: true})

	cmdRename := parser.NewCommand("r", "(Use: beryl r -i id -n newProjecName) Rename the selected project (ID can be viewed in --showall)")
	getId := cmdRename.Int("i", "id", &argparse.Options{Required: true})
	getRename := cmdRename.String("n", "name", &argparse.Options{Required: true})

	cmdReplace := parser.NewCommand("rp", "(Use: beryl rp -n projectName -w newProjectLocation) Changes in the internal db map to the project folder. (THIS DOES NOT REPLACE FILES OR FOLDERS) - -e|--verbose as optional")
	getReplace := cmdReplace.String("n", "name", &argparse.Options{Required: true})
	getRPNewFolder := cmdReplace.String("w", "newfolder", &argparse.Options{Required: true})
	getRPVerbose := cmdReplace.Flag("e", "verbose", &argparse.Options{Required: false})

	cmdDelete := parser.NewCommand("d", "(Use: beryl d -n projectName) Delete in the internal db map the project. (THIS DOES NOT DELETE FILES OR FOLDERS)")
	getDelete := cmdDelete.String("n", "name", &argparse.Options{Required: true})

	cmdAbout := parser.NewCommand("about", `Shows the "About" text`)

	//forceps flags
	switch {
	case len(args) < 2:
		return parser.Usage(parser), nil

	case (args[1] == "-v") || (args[1] == "--version"):
		return utils.ProgramName + " - " + utils.Version + "", nil

	case (args[1] == "-h") || (args[1] == "--help"):
		return parser.Usage(parser), nil
	}

	//commands
	err := parser.Parse(args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		return parser.Usage(err), nil
	}

	switch {
	case *cmdVersion:
		return utils.ProgramName + " - " + utils.Version, nil

	case cmdShowAll.Happened():
		return functional.ListProjectData("")

	case cmdShow.Happened():
		return functional.ListProjectData(*getShow)

	case cmdVerifyAll.Happened():
		result, err := functional.VerifyProjects(*getVAVerbose)
		return strings.Join(result, "\n "), err

	case cmdVerify.Happened():
		result, err := functional.VerifyAProject(*getVerify, *getVRVerbose)
		return strings.Join(result, "\n "), err

	case cmdAddNew.Happened():
		result, err := functional.AddProject(*getANName, *getANLocation, *getANVerbose)
		return strings.Join(result, "\n "), err

	case cmdAddHere.Happened():
		location, err := os.Getwd()
		if err != nil {
			return "", err
		}
		result, err := functional.AddProject(*getAHName, location, *getAHVerbose)
		return strings.Join(result, "\n "), err

	case cmdUpAll.Happened():
		result, err := functional.UpdateProjects(*getUAVerbose, *getUAForce)
		return strings.Join(result, "\n "), err

	case cmdUpdate.Happened():
		result, err := functional.UpdateAProject(*getUpdate, *getUVerbose, *getUForce)
		return strings.Join(result, "\n "), err

	case cmdTest.Happened():
		result, err := functional.TestAConnection(*getTest)
		return strings.Join(result, "\n "), err

	case cmdRename.Happened():
		result, err := functional.RenameAProject(*getId, *getRename)
		return strings.Join(result, "\n "), err

	case cmdReplace.Happened():
		result, err := functional.ReplaceAProject(*getReplace, *getRPNewFolder, *getRPVerbose)
		return strings.Join(result, "\n "), err

	case cmdDelete.Happened():
		result, err := functional.DeleteAProject(*getDelete)
		return strings.Join(result, "\n "), err

	case cmdAbout.Happened():
		return utils.AboutText, nil
	}

	return parser.Usage(err), nil
}
