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

const testDirectory = "testdata"

func Test(t *testing.T, solver func(in *bufio.Reader, out *bufio.Writer)) {
	files, err := os.ReadDir(testDirectory)
	if err != nil {
		panic(fmt.Errorf("reading directory: %w", err))
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".a") {
			continue
		}

		TestFile(t, file.Name(), solver)
	}
}

func TestWithLineComparator(t *testing.T, solver func(in *bufio.Reader, out *bufio.Writer), lineComparator func(expected, actual string) bool) {
	files, err := os.ReadDir(testDirectory)
	if err != nil {
		panic(fmt.Errorf("reading directory: %w", err))
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".a") {
			continue
		}

		TestFileWithLineComparator(t, file.Name(), solver, lineComparator)
	}
}

func TestFile(t *testing.T, file string, solver func(in *bufio.Reader, out *bufio.Writer)) {
	t.Run(file, func(t *testing.T) {
		t.Parallel()

		result := cleanResult(getResult(path.Join(testDirectory, file), solver))
		expected := cleanResult(getExpected(path.Join(testDirectory, file) + ".a"))

		assert.Equal(t, expected, result, fmt.Sprintf("testing file %s", file))
	})
}

func TestFileWithLineComparator(t *testing.T, file string, solver func(in *bufio.Reader, out *bufio.Writer), lineComparator func(expected, actual string) bool) {
	t.Run(file, func(t *testing.T) {
		t.Parallel()

		fmt.Printf("testing file %s\n", file)

		result := cleanResult(getResult(path.Join(testDirectory, file), solver))
		expected := cleanResult(getExpected(path.Join(testDirectory, file) + ".a"))

		resultLines := strings.Split(result, "\n")
		expectedLines := strings.Split(expected, "\n")

		assert.Equal(t, len(expectedLines), len(resultLines), fmt.Sprintf("testing file %s, result line count", file))

		for i := 0; i < len(resultLines); i++ {
			assert.True(
				t,
				lineComparator(expectedLines[i], resultLines[i]),
				fmt.Sprintf("testing file %s, line %s == %s", file, expectedLines[i], resultLines[i]),
			)
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
