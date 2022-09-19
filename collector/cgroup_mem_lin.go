package collector

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	reParens = regexp.MustCompile(`\((.*)\)`)
)

func (c *cgroupmemCollector) getCgroupUsageMem() (float64, error) {
	file, err := os.Open(cgroupFilePath("memory/memory.usage_in_bytes"))
	if err != nil {
		return -1, err
	}
	defer file.Close()

	return parseCgroupMem(file)
}

func parseCgroupMem(r io.Reader) (float64, error) {
	content, err := ioutil.ReadAll(r)
	contentf64, err := strconv.ParseFloat(strings.Replace(string(content), "\n", "", -1), 64)
	if err != nil {
		return -1, fmt.Errorf("invalid value in cgroupmemusage: %w", err)
	}
	return contentf64, err

}
