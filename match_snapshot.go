package snapshotmacther

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/onsi/gomega/matchers"
)

var snapshotDir = "__snapshots__"

func UpdateSnapshot() {
	_ = os.Setenv("UPDATE_SNAPSHOT", "1")
}

func MatchSnapshot(names ...string) *SnapshotMatcher {
	return &SnapshotMatcher{
		filename:    filepath.Join(names...),
		forceUpdate: os.Getenv("UPDATE_SNAPSHOT") != "",
	}
}

type SnapshotMatcher struct {
	filename    string
	forceUpdate bool
	matchers.EqualMatcher
}

func (matcher *SnapshotMatcher) Match(actual interface{}) (success bool, err error) {
	f := filepath.Join(snapshotDir, matcher.filename)

	var input string

	switch v := actual.(type) {
	case []byte:
		input = string(v)
	case string:
		input = v
	default:
		return false, fmt.Errorf("snapshot not support %T", actual)
	}

	if matcher.forceUpdate {
		if err := snapshot(f, []byte(input)); err != nil {
			return false, err
		}
		return true, nil
	}

	data, err := ioutil.ReadFile(f)
	if err != nil {
		if os.IsNotExist(err) {
			if err := snapshot(f, []byte(input)); err != nil {
				return false, err
			}
			return true, nil
		}
		return false, err
	}

	matcher.EqualMatcher.Expected = string(data)
	return matcher.EqualMatcher.Match(input)
}

func snapshot(filename string, input []byte) error {
	dirname := filepath.Dir(filename)
	if err := os.MkdirAll(dirname, os.ModePerm); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filename, input, os.ModePerm); err != nil {
		return err
	}
	return nil
}
