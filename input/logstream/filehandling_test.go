package logstream

import (
	gs "github.com/rafrombrc/gospec/src/gospec"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

func FilehandlingSpec(c gs.Context) {
	here, err := os.Getwd()
	c.Assume(err, gs.IsNil)
	dirPath := filepath.Join(here, "testdir")

	c.Specify("The directory scanner", func() {

		c.Specify("scans a directory properly", func() {
			matchRegex := regexp.MustCompile(dirPath + `/subdir/.*\.log(\..*)?`)
			results := ScanDirectoryForLogfiles(dirPath, matchRegex)
			c.Expect(len(results), gs.Equals, 3)
		})

		c.Specify("scana a directory with a bad regexp", func() {
			matchRegex := regexp.MustCompile(dirPath + "/subdir/.*.logg(.*)?")
			results := ScanDirectoryForLogfiles(dirPath, matchRegex)
			c.Expect(len(results), gs.Equals, 0)
		})
	})

	c.Specify("Populating logfile with match parts", func() {
		logfile := Logfile{}

		c.Specify("works without errors", func() {
			subexpNames := []string{"MonthName", "LogNumber"}
			matches := []string{"October", "24"}
			translation := make(SubmatchTranslationMap)
			logfile.PopulateMatchParts(subexpNames, matches, translation)
			c.Expect(logfile.MatchParts["MonthName"], gs.Equals, 10)
			c.Expect(logfile.MatchParts["LogNumber"], gs.Equals, 24)
		})

		c.Specify("works with bad month name", func() {
			subexpNames := []string{"MonthName", "LogNumber"}
			matches := []string{"Octoberrr", "24"}
			translation := make(SubmatchTranslationMap)
			err := logfile.PopulateMatchParts(subexpNames, matches, translation)
			c.Assume(err, gs.Not(gs.IsNil))
			c.Expect(err.Error(), gs.Equals, "Unable to locate month name: Octoberrr")
		})

		c.Specify("works with missing value in submatch translation map", func() {
			subexpNames := []string{"MonthName", "LogNumber"}
			matches := []string{"October", "24"}
			translation := make(SubmatchTranslationMap)
			translation["LogNumber"] = make(MatchTranslationMap)
			translation["LogNumber"]["23"] = 22
			err := logfile.PopulateMatchParts(subexpNames, matches, translation)
			c.Assume(err, gs.Not(gs.IsNil))
			c.Expect(err.Error(), gs.Equals, "Unable to locate value: (24) in translation map: LogNumber")
		})

		c.Specify("works with custom value in submatch translation map", func() {
			subexpNames := []string{"MonthName", "LogNumber"}
			matches := []string{"October", "24"}
			translation := make(SubmatchTranslationMap)
			translation["LogNumber"] = make(MatchTranslationMap)
			translation["LogNumber"]["24"] = 2
			logfile.PopulateMatchParts(subexpNames, matches, translation)
			c.Expect(logfile.MatchParts["MonthName"], gs.Equals, 10)
			c.Expect(logfile.MatchParts["LogNumber"], gs.Equals, 2)
		})
	})

	c.Specify("Populating logfiles with match parts", func() {
		translation := make(SubmatchTranslationMap)
		matchRegex := regexp.MustCompile(dirPath + `/subdir/.*\.log\.?(?P<FileNumber>.*)?`)
		logfiles := ScanDirectoryForLogfiles(dirPath, matchRegex)

		c.Specify("is populated", func() {
			logfiles.PopulateMatchParts(matchRegex, translation)
			c.Expect(len(logfiles), gs.Equals, 3)
			c.Expect(logfiles[0].MatchParts["FileNumber"], gs.Equals, -1)
			c.Expect(logfiles[1].MatchParts["FileNumber"], gs.Equals, 1)
		})

		c.Specify("returns errors", func() {
			translation["FileNumber"] = make(MatchTranslationMap)
			translation["FileNumber"]["23"] = 22
			err := logfiles.PopulateMatchParts(matchRegex, translation)
			c.Assume(err, gs.Not(gs.IsNil))
			c.Expect(len(logfiles), gs.Equals, 3)
		})
	})

	c.Specify("Sorting logfiles", func() {
		translation := make(SubmatchTranslationMap)
		matchRegex := regexp.MustCompile(dirPath + `/subdir/.*\.log\.?(?P<FileNumber>.*)?`)
		logfiles := ScanDirectoryForLogfiles(dirPath, matchRegex)
		err := logfiles.PopulateMatchParts(matchRegex, translation)
		c.Assume(err, gs.IsNil)
		c.Expect(len(logfiles), gs.Equals, 3)

		c.Specify("can be sorted newest to oldest", func() {
			byp := ByPriority{Logfiles: logfiles, Priority: []string{"FileNumber"}}
			sort.Sort(byp)
			c.Expect(logfiles[0].MatchParts["FileNumber"], gs.Equals, -1)
			c.Expect(logfiles[1].MatchParts["FileNumber"], gs.Equals, 1)
		})

		c.Specify("can be sorted oldest to newest", func() {
			byp := ByPriority{Logfiles: logfiles, Priority: []string{"^FileNumber"}}
			sort.Sort(byp)
			c.Expect(logfiles[0].MatchParts["FileNumber"], gs.Equals, 2)
			c.Expect(logfiles[1].MatchParts["FileNumber"], gs.Equals, 1)
		})
	})
}
