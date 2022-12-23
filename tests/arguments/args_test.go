package tests

import (
	"log"
	"os"
	"testing"

	"github.com/mrGlasses/BerylSQLHelper/arguments"
	"github.com/mrGlasses/BerylSQLHelper/utils"
)

func TestExecuteArgumentsVersion(t *testing.T) {
	var args []string
	local, _ := os.Executable()
	args = append(args, local, "-v")
	result, err := arguments.ExecuteArguments(args)
	if result != (utils.ProgramName + " - " + utils.Version) {
		log.Fatalf("TestExecuteArgumentsVersion failed \n Result: %v \n Error: %v", result, err)
		log.Fatalln("" + result)
	}
}

func TestExecuteArgumentsAbout(t *testing.T) {
	var args []string
	local, _ := os.Executable()
	args = append(args, local, "about")
	result, err := arguments.ExecuteArguments(args)
	if result != utils.AboutText {
		log.Fatalf("TestExecuteArgumentsAbout failed \n Result: %v \n Error: %v", result, err)
		log.Fatalln("" + result)
	}
}
