package pluginmanager_test

import (
	"github.com/fredwangwang/bosh-plugin/pluginmanager"
	. "github.com/onsi/gomega"
	"testing"
)

func TestValidatePluginName(t *testing.T) {
	g := NewGomegaWithT(t)

	g.Expect(pluginmanager.ValidatePluginName("a2C-_D")).To(BeNil())
	g.Expect(pluginmanager.ValidatePluginName("a+b=c")).To(HaveOccurred())
	g.Expect(pluginmanager.ValidatePluginName("invalid with space")).To(HaveOccurred())
}
