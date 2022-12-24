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

	// args = append(args, "sa") //test arguments
	fmt.Println(args)

	parser := argparse.NewParser(utils.CommandName, utils.ProgramDescription)

	cmdVersion := parser.Flag("v", "version", &argparse.Options{Required: false, Help: "Shows the installed version of the code"})

	cmdShowAll := parser.NewCommand("sa", "Shows all main folders for each project")

	cmdShow := parser.NewCommand("s", "(Use: s -n projectName) Shows the data of the selected project")
	getShow := cmdShow.String("n", "name", &argparse.Options{Required: true})

	cmdVerifyAll := parser.NewCommand("va", "Verifies all projects and covered folders for updates - -e|--verbose as optional")
	getVAVerbose := cmdVerifyAll.Flag("e", "verbose", &argparse.Options{Required: false})

	cmdVerify := parser.NewCommand("vr", "(Use: vr -n projectName) Verifies a specific project and covered folders for updates - -e|--verbose as optional")
	getVerify := cmdVerify.String("n", "name", &argparse.Options{Required: true})
	getVRVerbose := cmdVerify.Flag("e", "verbose", &argparse.Options{Required: false})

	cmdAddNew := parser.NewCommand("an", "(Use: an -n projectName -l projectLocation) Adds a new project and its folder to the app - -e|--verbose as optional")
	getANName := cmdAddNew.String("n", "name", &argparse.Options{Required: true})
	getANLocation := cmdAddNew.String("l", "location", &argparse.Options{Required: true})
	getANVerbose := cmdAddNew.Flag("e", "verbose", &argparse.Options{Required: false})

	cmdAddHere := parser.NewCommand("ah", "(Use: ah -n projectName) Adds the current folder to the app")
	getAHName := cmdAddHere.String("n", "name", &argparse.Options{Required: true})
	getAHVerbose := cmdAddHere.Flag("e", "verbose", &argparse.Options{Required: false})

	// cmdUpAll := parser.String("ua", "updateall", &argparse.Options{Required: false, Help: "Updates all projects added to the app", Default: ""})

	// cmdVUpdate := parser.String("u", "update", &argparse.Options{Required: false, Help: "(Use: -u projectName) Updates a specific project", Default: ""})

	// cmdForce := parser.String("f", "force", &argparse.Options{Required: false, Help: "(Use: -u projectName -f) Only works with -u command - (be careful) Re-run all files in all folders.", Default: ""})

	// cmdTest := parser.String("tc", "testconnection", &argparse.Options{Required: false, Help: "(Use: -tc projectName) Test the connection with the server/database", Default: ""})

	// cmdRename := parser.Int("r", "rename", &argparse.Options{Required: false, Help: "(Use: -r id) Rename the selected project (ID can be viewed in --showall)", Default: ""})

	// cmdReplace := parser.String("rp", "replace", &argparse.Options{Required: false, Help: "(Use: -rp projectName -rp newProjectLocation) Changes in the internal db map to the project folder. (THIS DOES NOT REPLACE FILES OR FOLDERS)", Default: ""})

	// cmdDelete := parser.String("del", "delete", &argparse.Options{Required: false, Help: "(Use: -del projectName) Delete in the internal db map the project. (THIS DOES NOT DELETE FILES OR FOLDERS)", Default: ""})

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

	case cmdAbout.Happened():
		return utils.AboutText, nil
	}

	return parser.Usage(err), nil
}
