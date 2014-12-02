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

	sourceTemplate = []byte(`package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, playground")
}
`)
)

type Config struct {
	List bool
	Edit bool
}

func playgroundRootPath() (string, error) {
	root, err := gitconfig.Global("goplay.root")
	if err == nil {
		return root, nil
	}

	if err != gitconfig.ErrNotFound {
		return "", err
	}

	return "/tmp/goplay", nil
}

func listPlaygroundDirs() error {
	root, err := playgroundRootPath()
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

func createPlaygroundDir() (string, error) {
	root, err := playgroundRootPath()
	if err != nil {
		return "", err
	}

	name := time.Now().Format("2006-01-02_15-04-05")
	path := filepath.Join(root, name)
	path, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(path, 0700); err != nil {
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

func getenv(name, defaultValue string) string {
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
			continue
		}
		return v, nil
	}
	return "", fmt.Errorf("editor not found, please set $EDITOR environment variable")
}

func createGoFile(dirPath string) (string, error) {
	path, err := filepath.Abs(filepath.Join(dirPath, "main.go"))
	if err != nil {
		return "", err
	}

	if err := ioutil.WriteFile(path, sourceTemplate, 0600); err != nil {
		return "", err
	}
	return path, err
}

func gotoPlayground(dirPath, filePath string) error {
	shell, err := exec.LookPath(getenv("SHELL", "/bin/sh"))
	if err != nil {
		return err
	}

	if err := os.Chdir(dirPath); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "\tcd %s\n", dirPath)

	var argv []string
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

	if err := os.Setenv("GOPATH", libs+":"+os.Getenv("GOPATH")); err != nil {
		return err
	}

	return syscall.Exec(shell, argv, syscall.Environ())
}

func main() {
	abortOnError := func(err error) {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	err := parseConfig()
	abortOnError(err)

	if config.List {
		err = listPlaygroundDirs()
		abortOnError(err)
		os.Exit(0)
	}

	dirPath, err := createPlaygroundDir()
	abortOnError(err)

	filePath, err := createGoFile(dirPath)
	abortOnError(err)

	err = gotoPlayground(dirPath, filePath)
	abortOnError(err)

	os.Exit(0)
}
