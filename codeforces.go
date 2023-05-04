package codeforces

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"strings"
	"testing"
)

type Tester struct {
	t                *testing.T
	solver           func(in *bufio.Reader, out *bufio.Writer)
	lineComparator   func(expected, actual string) bool
	outputNormalizer func(output string) string
	testDataDir      string
}

func New(t *testing.T, solver func(in *bufio.Reader, out *bufio.Writer)) *Tester {
	return &Tester{
		t:           t,
		solver:      solver,
		testDataDir: "testdata",
	}
}

func (tester *Tester) WithLineComparator(lineComparator func(expected, actual string) bool) *Tester {
	tester.lineComparator = lineComparator
	return tester
}

func (tester *Tester) WithOutputNormalizer(outputNormalizer func(output string) string) *Tester {
	tester.outputNormalizer = outputNormalizer
	return tester
}

func (tester *Tester) Test() {
	files, err := os.ReadDir(tester.testDataDir)
	if err != nil {
		panic(fmt.Errorf("reading directory: %w", err))
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".a") {
			continue
		}

		tester.TestFile(file.Name())
	}
}

func (tester *Tester) TestFile(file string) {
	tester.t.Run(file, func(t *testing.T) {
		t.Parallel()

		fmt.Printf("testing file %s\n", file)

		result := cleanResult(getResult(path.Join(tester.testDataDir, file), tester.solver))
		if tester.outputNormalizer != nil {
			result = tester.outputNormalizer(result)
		}

		expected := cleanResult(getExpected(path.Join(tester.testDataDir, file) + ".a"))
		if tester.outputNormalizer != nil {
			expected = tester.outputNormalizer(expected)
		}

		resultLines := strings.Split(result, "\n")
		expectedLines := strings.Split(expected, "\n")

		assert.Equal(t, len(expectedLines), len(resultLines), "result line count does not match")

		for i := 0; i < len(resultLines); i++ {
			if tester.lineComparator != nil {
				assert.True(
					t,
					tester.lineComparator(expectedLines[i], resultLines[i]),
					fmt.Sprintf("line comparator failed for line %d: line from expected is %q, line from result is %q", i, expectedLines[i], resultLines[i]),
				)
			} else {
				assert.Equal(t, expected, result, fmt.Sprintf("testing file %s", file))
			}
		}
	})
}

func getExpected(file string) string {
	data, err := os.ReadFile(file)
	if err != nil {
		panic(fmt.Errorf("reading test result file: %w", err))
	}

	return string(data)
}

func getResult(file string, solver func(in *bufio.Reader, out *bufio.Writer)) string {
	inFile, err := os.Open(file)
	if err != nil {
		panic(fmt.Errorf("opening file: %w", err))
	}
	defer inFile.Close()

	buf := bytes.Buffer{}
	out := bufio.NewWriter(&buf)

	solver(bufio.NewReader(inFile), out)

	out.Flush()
	return buf.String()
}

func cleanResult(text string) string {
	text = strings.Replace(text, "\r\n", "\n", -1)
	text = strings.TrimRight(text, "\n")
	return text
}
