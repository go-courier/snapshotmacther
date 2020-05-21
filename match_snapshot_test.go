package snapshotmacther

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestSnapshot(t *testing.T) {
	scope := t.Name()

	t.Run("match snapshot", func(t *testing.T) {
		gomega.NewWithT(t).Expect("111").To(MatchSnapshot(scope, "test.txt"))
	})

	t.Run("not match snapshot", func(t *testing.T) {
		gomega.NewWithT(t).Expect("222").NotTo(MatchSnapshot(scope, "test.txt"))
	})

	t.Run("update snapshot", func(t *testing.T) {
		UpdateSnapshot()
		gomega.NewWithT(t).Expect("222").To(MatchSnapshot(scope, "test.txt"))
		gomega.NewWithT(t).Expect("111").To(MatchSnapshot(scope, "test.txt"))
	})
}
