package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	var countPtr = flag.String("count", "10", "number of commits to walk back on")
	flag.Parse()

	lastMergeDurations(*countPtr)
}

func lastMergeDurations(count string) {
	mergesLines := getMergeCommits(count)

	for _, line := range mergesLines {
		mergeDuration(line)
	}
}

func mergeDuration(line mergesLine) {
	branchTime := getBranchTime(line.mergeSha, line.mergeBaseSha)

	mergeTimeInt, _ := strconv.ParseInt(line.mergeTime, 10, 64)
	branchTimeInt, _ := strconv.ParseInt(branchTime, 10, 64)

	t1 := time.Unix(mergeTimeInt, 0)
	t2 := time.Unix(branchTimeInt, 0)
	delta := t1.Sub(t2)
	fmt.Println(delta, line.mergeSha, line.subject)
}

// get oldest commit unix timestamp
// git rev-list --pretty='%at' 5d8f7fca0...b439a2264 | tail -1
// commit dc17abef73365b9f659c5cb2655fc59404720340
// 1539970561
// commit 2a0b532200f2751daceb1af6b1e55285cbb836af
// 1539885622
func getBranchTime(mergeCommit, mergeBase string) string {
	commitRange := mergeCommit + "..." + mergeBase
	cmdArgs := []string{"rev-list", "--pretty=%at", commitRange}

	lines := runGitCmd(cmdArgs)
	lastDate := lines[len(lines)-1]

	return lastDate
}

type mergesLine struct {
	subject         string
	mergeSha        string
	mergeBaseSha    string
	latestCommitSha string
	mergeTime       string
}

// [subject, merge-sha, merge-base-sha, latest-commit-sha, merge-time]
func getMergeCommits(count string) []mergesLine {
	numberCommits := "-n " + count
	cmdArgs := []string{"log", "--merges", "--pretty=%f %H %P %at", numberCommits}
	lines := runGitCmd(cmdArgs)

	mergesLines := []mergesLine{}
	for _, line := range lines {
		values := strings.Split(line, " ")
		mergesLines = append(mergesLines, mergesLine{
			subject:         string(values[0]),
			mergeSha:        string(values[1]),
			mergeBaseSha:    string(values[2]),
			latestCommitSha: string(values[3]),
			mergeTime:       string(values[4]),
		})
	}

	return mergesLines
}

func runGitCmd(cmdArgs []string) []string {
	cmdOut, err := exec.Command("git", cmdArgs...).Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error running cmd: ", cmdArgs, err)
		os.Exit(1)
	}

	resultString := string(cmdOut[:])
	resultLines := strings.Split(resultString, "\n")
	if len(resultLines) > 0 {
		resultLines = resultLines[:len(resultLines)-1]
	}

	return resultLines
}
