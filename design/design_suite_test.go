package design

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDesign(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Design Suite")
}
