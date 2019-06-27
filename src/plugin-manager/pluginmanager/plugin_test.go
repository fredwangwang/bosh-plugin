package pluginmanager_test

import (
	"github.com/fredwangwang/bosh-plugin-manager/pluginmanager"
	. "github.com/onsi/gomega"
	"testing"
)

func TestValidatePluginName(t *testing.T) {
	g := NewGomegaWithT(t)

	g.Expect(pluginmanager.ValidatePluginName("a2C-_ D")).To(BeNil())
	g.Expect(pluginmanager.ValidatePluginName("a+b=c")).To(HaveOccurred())
}
