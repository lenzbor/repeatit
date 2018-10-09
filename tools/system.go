package tools

import (
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// RunCmdWithErrorCode executes a command and returns the output of the command.
func RunCmdWithErrorCode(checkoutFolder string, cmd exec.Cmd, errMsg string, debug, showCmd bool) (string, int, error) {
	if _, err := os.Stat(checkoutFolder); os.IsNotExist(err) {
		return "", 0, fmt.Errorf("the folder %q does not exist", checkoutFolder)
	}
	cmd.Dir = checkoutFolder
	if showCmd || debug {
		PositiveStatus(fmt.Sprintf("%s", strings.Join(cmd.Args, " ")))
	}

	var stdout, stderr []byte
	var err error
	var errorcode int
	if debug {
		stdout, stderr, errorcode, err = ExecCmdCaptureStreamsPublishOutput(cmd)
	} else {
		stdout, stderr, errorcode, err = ExecCmdCaptureStreamsNoOutput(cmd)
	}
	output := strings.Join([]string{string(stdout), string(stderr)}, "\n")
	if err != nil {
		return "", errorcode, errors.Wrap(err, fmt.Sprintf("%s due to %q", errMsg, output))
	}
	return output, errorcode, nil
}

// RunCmd is the same as RunCmdWithErrorCode but no error code from the command
// is returned.
func RunCmd(checkoutFolder string, cmd exec.Cmd, errMsg string, debug, showCmd bool) (string, error) {
	output, _, err := RunCmdWithErrorCode(checkoutFolder, cmd, errMsg, debug, showCmd)
	return output, err
}

// CreateSymlink executes the ln command.
func CreateSymlink(folder, source, dest string, debug, showCmd bool) (string, error) {
	cmd := exec.Command("ln", "-sf", source, dest)
	errorMsg := fmt.Sprintf("error while trying to create symlink from %s to %s in %s", source, dest, folder)
	output, err := RunCmd(folder, *cmd, errorMsg, debug, showCmd)

	return output, err
}

// FileExists is a function to improve code readibility. It tells if a file
// exists or not.
func FileExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return !fi.IsDir(), nil
}

// DirExists is a function to improve code readibility. It tells if a directory
// exists or not.
func DirExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return fi.IsDir(), nil
}

// Copy creates a copy of the source denoted by its absolute path 'sourceItem'
// and creates a copy as the absolute path 'targetItem'. This function works
// for both files and folders.
// It requires that rsync is in your PATH.
func Copy(targetItem, sourceItem string) error {
	cmd := exec.Command("rsync", "-a", sourceItem, targetItem)
	var stdout, stderr []byte
	var err error
	if IsDebugActivated() {
		stdout, stderr, _, err = ExecCmdCaptureStreamsPublishOutput(*cmd)
	} else {
		stdout, stderr, _, err = ExecCmdCaptureStreamsNoOutput(*cmd)
	}
	if err != nil {
		output := strings.Join([]string{string(stdout), string(stderr)}, "\n")
		return errors.Wrap(err, fmt.Sprintf("error while copying file %q to %q. (%s)", sourceItem, targetItem, output))
	}
	return nil
}

// BuildSymlinksChain builds the full list of links that point to
// each other. Maximum level is 255 (unix limit ?)
// This code is *heavily* inspired from the function
// walkLink from go source path/filepath/symlink.go
func BuildSymlinksChain(path string, symlinksChain *[]string) ([]string, error) {
	if len(*symlinksChain) == 0 {
		*symlinksChain = append(*symlinksChain, path)
	}
	if len(*symlinksChain) > 255 {
		return []string{}, fmt.Errorf("too many symlinks chained (>255)")
	}
	fi, err := os.Lstat(path)
	if err != nil {
		return *symlinksChain, err
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		return *symlinksChain, nil
	}
	newpath, err := os.Readlink(path)
	if err != nil {
		return *symlinksChain, err
	}
	if !filepath.IsAbs(newpath) {
		newpath = filepath.Clean(filepath.Join(filepath.Dir(path), newpath))
	}
	*symlinksChain = append(*symlinksChain, newpath)

	return BuildSymlinksChain(newpath, symlinksChain)
}

// FindFile searchs recursively a file in a folder (search is case
// sensitive).
func FindFile(filename, folder string) ([]string, error) {
	return findItemInFolder(filename, "f", folder)
}

// FindFolder searchs recursively a folder in a folder (search is case
// sensitive).
func FindFolder(filename, folder string) ([]string, error) {
	return findItemInFolder(filename, "d", folder)
}

// relies on the find command
func findItemInFolder(itemName, itemtype, folder string) ([]string, error) {
	output := []string{}
	cmd := exec.Command("find", ".", "-type", itemtype, "-name", itemName)
	cmd.Dir = folder
	var stdout, stderr []byte
	var err error
	if IsDebugActivated() {
		stdout, stderr, _, err = ExecCmdCaptureStreamsPublishOutput(*cmd)
	} else {
		stdout, stderr, _, err = ExecCmdCaptureStreamsNoOutput(*cmd)
	}
	if err != nil {
		return output, errors.Wrap(err, fmt.Sprintf("error while searching %s in %s (%s)", itemName, folder, string(stderr)))
	}
	out := strings.Trim(string(stdout), "\n")
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.Index(line, "/vendor") == -1 && strings.Trim(line, " ") != "" {
			Debug(fmt.Sprintf("Adding %s", filepath.Clean(strings.Replace(line, "./", folder+"/", -1))))
			output = append(output, filepath.Clean(strings.Replace(line, "./", folder+"/", -1)))
		}
	}

	return output, nil
}

// WriteASTToFile saves the AST structure to a GoFile.
func WriteASTToFile(pathToGoFile string, fset *token.FileSet, astFile *ast.File) error {
	// Rewriting the file. For safety, creating a file with suffix
	pathToTmpFile := pathToGoFile + ".migrated"
	migrated, err := os.Create(pathToTmpFile)
	if err != nil {
		return errors.Wrap(err, "failed to create a file to host imports migration.")
	}
	// write changes to .temp file, and include proper formatting.
	err = (&printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}).Fprint(migrated, fset, astFile)
	if err != nil {
		return err
	}
	// close the writer
	err = migrated.Close()
	if err != nil {
		return err
	}
	err = os.Rename(pathToTmpFile, pathToGoFile)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to rename migrated file %q to original file %q", pathToTmpFile, pathToGoFile))
	}
	time.Sleep(time.Millisecond * 10)

	return nil
}

// WriteASTToFs writes an AST to an abstract file system.
func WriteASTToFs(fs afero.Fs, pathToGoFile string, fset *token.FileSet, astFile *ast.File) error {
	f, err := fs.Create(pathToGoFile)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to create file %q on afero fs", pathToGoFile))
	}
	// write changes to .temp file, and include proper formatting.
	err = (&printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}).Fprint(f, fset, astFile)
	if err != nil {
		return err
	}
	// close the writer
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}

// ListFolders lists the folders in a given tree structure. The returned list
// is sorted.
func ListFolders(path string, ignorelist map[string]int) ([]string, error) {
	d, err := os.Open(path)
	if err != nil {
		return []string{}, errors.Wrap(err, fmt.Sprintf("error while listing content of %s", path))
	}
	fis, err := d.Readdir(-1)
	if err != nil {
		return []string{}, errors.Wrap(err, fmt.Sprintf("failed to list content of folder %s", path))
	}
	list := []string{}
	for _, fi := range fis {
		if !fi.IsDir() {
			continue
		}
		if fi.Name() == ".git" {
			continue
		}
		// Now we are sure that the file info is a directory and is not .git
		if ignorelist[fi.Name()] != 0 {
			continue
		}
		list = append(list, strings.Replace(filepath.Join(path, fi.Name()), filepath.Join(os.Getenv("GOROOT"), "src")+"/", "", -1))
		subdir := filepath.Join(path, fi.Name())
		sublist, err := ListFolders(subdir, ignorelist)
		if err != nil {
			return []string{}, errors.Wrap(err, fmt.Sprintf("failed to list content of %s", subdir))
		}
		list = append(list, sublist...)
	}
	d.Close()
	sort.Strings(list)
	return list, nil
}

// CheckCommandIsAvailable tells if a command is available in the PATH of the user.
func CheckCommandIsAvailable(commandToCheck string) (bool, error) {
	cmd := exec.Command("command", "-v", commandToCheck)
	errorMsg := fmt.Sprintf("error while checking if command %s is in PATH", commandToCheck)
	// folder where the command must be is useless since it is not used.
	// So using "/"
	// false means : do not show the executed command to the end user.
	_, errorcode, err := RunCmdWithErrorCode("/", *cmd, errorMsg, IsDebugActivated(), false)
	if err != nil {
		return false, errors.Wrap(err, fmt.Sprintf("error while checking if command %s is present in user's PATH", commandToCheck))
	}

	if errorcode == 0 {
		return true, nil
	}
	return false, nil
}

// SaveBytesToFile writes an array of bytes to an output file. The
// function assures you that it creates the tree structure to the output
// file and report any error.
func SaveBytesToFile(content []byte, outputFile string, overwrite bool) error {
	exists, err := FileExists(outputFile)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error while checking if the output file %q already exists", outputFile))
	}
	if exists && !overwrite {
		return fmt.Errorf("cannot persist the data since the output file %q already exists", outputFile)
	}
	parentFolder := filepath.Dir(outputFile)
	exists, err = DirExists(parentFolder)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error while checking if the parent folder of %q exists", outputFile))
	}
	if !exists {
		err = os.MkdirAll(parentFolder, DefaultPermission)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to create tree structure to %q to store the data", parentFolder))
		}
	}
	err = ioutil.WriteFile(outputFile, content, DefaultPermission)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error while saving the data to file %q", outputFile))
	}
	return nil
}
