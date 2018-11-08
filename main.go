package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// # get: merge-commit | merge-base | last commit | merge-commit-timestamp
// git log --merges --pretty='%H %P %at' -n 10

// # with: merge-commit...merge-base
// # get oldest commit unix timestamp
// git rev-list --pretty='%at' 5d8f7fca0...b439a2264 | tail -1
// 1539885622

func main() {
	commits := getMergeCommits()

	// fmt.Println(commits)

	for _, mcSring := range commits {
		body := strings.Split(mcSring, " ")
		mergeCommit := string(body[0])
		mergeBase := string(body[1])
		mergeTime := string(body[3])

		// fmt.Println(mergeCommit)
		// fmt.Println(mergeBase)
		branchTime := getBranchTime(mergeCommit, mergeBase)
		// fmt.Println(branchTime)

		// a, err := time.Parse(mergeTime, branchTime)
		// if err != nil {
		// 	return
		// }

		mergeTimeInt, _ := strconv.ParseInt(mergeTime, 10, 64)
		branchTimeInt, _ := strconv.ParseInt(branchTime, 10, 64)

		t1 := time.Unix(mergeTimeInt, 0)
		t2 := time.Unix(branchTimeInt, 0)
		delta := t1.Sub(t2)
		fmt.Println(delta)
	}
}

func getBranchTime(mergeCommit, mergeBase string) string {
	// 	with: merge-commit...merge-base
	// get oldest commit unix timestamp
	// git rev-list --pretty='%at' 5d8f7fca0...b439a2264 | tail -1
	// commit dc17abef73365b9f659c5cb2655fc59404720340
	// 1539970561
	// commit 2a0b532200f2751daceb1af6b1e55285cbb836af
	// 1539885622
	var (
		cmdOut []byte
		err    error
	)

	cmdName := "git"
	commitRange := mergeCommit + "..." + mergeBase
	cmdArgs := []string{"rev-list", "--pretty=%at", commitRange}
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running git rev-parse command: ", err)
		os.Exit(1)
	}

	stringResult := string(cmdOut[:])

	lines := strings.Split(stringResult, "\n")
	if len(lines) > 0 {
		lines = lines[:len(lines)-1]
	}

	lastDate := lines[len(lines)-1]

	return lastDate
}

func getMergeCommits() []string {
	var (
		cmdOut []byte
		err    error
	)
	cmdName := "git"
	cmdArgs := []string{"log", "--merges", "--pretty=%H %P %at", "-n 5"}
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running git rev-parse command: ", err)
		os.Exit(1)
	}
	stringResult := string(cmdOut[:])
	resultArgs := strings.Split(stringResult, "\n")
	if len(resultArgs) > 0 {
		resultArgs = resultArgs[:len(resultArgs)-1]
	}

	return resultArgs
}
