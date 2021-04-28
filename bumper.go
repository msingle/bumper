package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const exitFail = 1
const extension = ".js" // change this to flag

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFail)
	}
}

func run(args []string, stdout io.Writer) error {
	if len(args) < 1 {
		return errors.New("no files")
	}
	for _, arg := range args[1:] {
		oldName := arg
		fmt.Fprintf(stdout, "Looking for %s\n", oldName)

		newName, err := version(oldName, extension)
		if err != nil {
			return err
		}
		// find all files that reference oldName
		out, err := exec.Command("rg", "-l", oldName).Output()
		if err != nil {
			return err
		}
		oldRefs := strings.Split(string(out), "\n")
		for _, oldRef := range oldRefs {
			subStr := fmt.Sprintf("s|%s|%s|", oldName, newName)
			subOut, err := exec.Command("sed", "-i", subStr, oldRef).CombinedOutput()
			if err != nil {
				fmt.Fprintf(stdout, "%s\n", subOut)
				log.Fatal(err)
			}
		}

	}
	return nil
}

func version(fnameWithExt, extension string) (string, error) {
	fname := strings.TrimSuffix(fnameWithExt, extension)
	if fname == fnameWithExt {
		return "", fmt.Errorf("file %q is not of type %q", fname, fnameWithExt)
	}
	vspot := strings.LastIndex(fname, "_v")
	if vspot == -1 {
		return "", fmt.Errorf("didn't find '_v' in %q", fname)
	}
	chunks := strings.Split(fname, "_v")
	raw := chunks[len(chunks)-1]
	bumped, err := bump(raw)
	if err != nil {
		return "", err
	}
	chunks[len(chunks)-1] = bumped
	joined := strings.Join(chunks, "_v")
	return joined + extension, nil
}

func bump(raw string) (string, error) {
	vStrings := strings.Split(raw, ".")
	lastIdx := len(vStrings) - 1
	last, err := strconv.Atoi(vStrings[lastIdx])
	if err != nil {
		return "", err
	}
	last++
	vStrings[lastIdx] = strconv.Itoa(last)

	return strings.Join(vStrings, "."), nil
}
