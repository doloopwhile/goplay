package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"syscall"
	"time"

	"github.com/tcnksm/go-gitconfig"
)

var (
	config  Config
	version string
)

const (
	timeFormat = "2006-01-02_15-04-05"
)

var sourceTemplate = []byte(`package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, playground")
}
`)

type Config struct {
	List bool
	Edit bool
}

func abort(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func getCampRootPath() (string, error) {
	root, err := gitconfig.Global("goplay.root")
	if err == nil {
		return root, nil
	}

	if err != gitconfig.ErrNotFound {
		return "", err
	}

	return "/tmp/goplay", nil
}

func listCampDirs() error {
	root, err := getCampRootPath()
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	names := []string{}

	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(names)))
	for _, name := range names {
		abs, err := filepath.Abs(filepath.Join(root, name))
		if err != nil {
			return err
		}
		fmt.Println(abs)
	}
	return nil
}

func createCampDir() (string, error) {
	root, err := getCampRootPath()
	if err != nil {
		return "", err
	}

	name := time.Now().Format(timeFormat)
	path := filepath.Join(root, name)
	path, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(path, 0700)
	if err != nil {
		return "", err
	}

	return path, nil
}

func parseConfig() error {
	flag.BoolVar(&config.Edit, "e", false, "Open created script with editor")
	flag.BoolVar(&config.List, "l", false, "List park directries")
	v := flag.Bool("version", false, "print app version")
	flag.Parse()

	if *v {
		fmt.Printf("goplay: %s\n", version)
		os.Exit(0)
	}

	return nil
}

func getenvDefault(name, defaultValue string) string {
	v := os.Getenv(name)
	if v == "" {
		return defaultValue
	}
	return v
}

func getEditorCommand() (string, error) {
	if v := os.Getenv("EDITOR"); v != "" {
		return v, nil
	}

	if v := os.Getenv("VISUAL"); v != "" {
		return v, nil
	}

	// in Unix-like OS, try to find editors
	for _, v := range []string{"editor", "vi", "nano"} {
		_, err := exec.LookPath(v)
		if err != nil {
			return v, nil
		}
	}
	return "", fmt.Errorf("editor not found, please set $EDITOR environment variable")
}

func createGoFile(dirPath string) (string, error) {
	path, err := filepath.Abs(filepath.Join(dirPath, "main.go"))
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(path, sourceTemplate, 0600)
	if err != nil {
		return "", err
	}
	return path, err
}

func gotoCamp(dirPath, filePath string) error {
	var err error

	shell := getenvDefault("SHELL", "/bin/sh")
	shell, err = exec.LookPath(shell)
	if err != nil {
		return err
	}

	err = os.Chdir(dirPath)
	if err != nil {
		return err
	}

	var argv []string

	fmt.Fprintf(os.Stderr, "\tcd %s\n", dirPath)

	if config.Edit {
		editor, err := getEditorCommand()
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "\t%s %s\n", editor, filePath)
		argv = []string{
			shell,
			"-c",
			fmt.Sprintf("%s %s; exec %s", editor, filePath, shell),
		}
	} else {
		argv = []string{shell}
	}

	libs, err := filepath.Abs(filepath.Join(dirPath, "golibs"))
	if err != nil {
		return err
	}

	err = os.Setenv("GOPATH", libs+":"+os.Getenv("GOPATH"))
	if err != nil {
		return err
	}

	return syscall.Exec(shell, argv, syscall.Environ())
}

func main() {
	if err := parseConfig(); err != nil {
		abort(err)
	}

	if config.List {
		if err := listCampDirs(); err != nil {
			abort(err)
		}
		os.Exit(0)
	}

	dirPath, err := createCampDir()
	if err != nil {
		abort(err)
	}

	filePath, err := createGoFile(dirPath)
	if err != nil {
		abort(err)
	}

	if err := gotoCamp(dirPath, filePath); err != nil {
		abort(err)
	}

	os.Exit(0)
}
