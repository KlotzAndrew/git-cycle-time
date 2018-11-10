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
	mergeCommits := getMergeCommits(count)
	avg := avgDuration(mergeCommits)

	for _, commit := range mergeCommits {
		fmt.Println(commit.duration, commit.mergeSha, commit.subject)
	}

	days := (avg * time.Second).Hours() / 24
	if days < float64(1) {
		fmt.Println("avg: ", time.Duration(avg*time.Second))
	} else {
		fmt.Println("avg: ", strconv.FormatFloat(days, 'f', 1, 64)+" days")
	}
}

func avgDuration(commits []mergeCommit) time.Duration {
	total := 0 * time.Second
	var count float64
	for _, commit := range commits {
		count++
		total += commit.duration
	}

	return time.Duration(total.Seconds() / count)
}

func getBranchTime(mergeCommit, mergeBase string) string {
	commitRange := mergeCommit + "..." + mergeBase
	cmdArgs := []string{"rev-list", "--pretty=%at", commitRange}

	lines := runGitCmd(cmdArgs)
	lastDate := lines[len(lines)-1]

	return lastDate
}

type mergeCommit struct {
	subject         string
	mergeSha        string
	mergeBaseSha    string
	latestCommitSha string
	mergeTime       string
	duration        time.Duration
}

func (c *mergeCommit) calcMergeDuration() {
	branchTime := getBranchTime(c.mergeSha, c.mergeBaseSha)

	mergeTimeInt, _ := strconv.ParseInt(c.mergeTime, 10, 64)
	branchTimeInt, _ := strconv.ParseInt(branchTime, 10, 64)

	t1 := time.Unix(mergeTimeInt, 0)
	t2 := time.Unix(branchTimeInt, 0)
	c.duration = t1.Sub(t2)
}

// [subject, merge-sha, merge-base-sha, latest-commit-sha, merge-time]
func getMergeCommits(count string) []mergeCommit {
	numberCommits := "-n " + count
	cmdArgs := []string{"log", "--merges", "--pretty=%f %H %P %at", numberCommits}
	lines := runGitCmd(cmdArgs)

	mergeCommits := []mergeCommit{}
	for _, line := range lines {
		commit := commitFromLine(line)
		commit.calcMergeDuration()
		mergeCommits = append(mergeCommits, commit)
	}

	return mergeCommits
}

func commitFromLine(line string) mergeCommit {
	values := strings.Split(line, " ")
	return mergeCommit{
		subject:         string(values[0]),
		mergeSha:        string(values[1]),
		mergeBaseSha:    string(values[2]),
		latestCommitSha: string(values[3]),
		mergeTime:       string(values[4]),
	}
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
