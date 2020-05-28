package conformance

/*<speck>
## Conformance

### Minimum Requirements

For a registry to be considered fully conformant against this specification, it must implement the HTTP endpoints
required by each of the four (4) major workflow categories:

1. **Pull** (REQUIRED) - Ability to fetch content from a registry
2. **Push** - Ability to publish content to a registry
3. **Content Discovery** - Ability to list or otherwise query the content stored in a registry
4. **Content Management** - Ability to delete (or otherwise manipulate) content stored in a registry

At a bare minimum, registries claiming to be "OCI-Compliant" MUST support all facets of the pull workflow.

In order to test a registry's conformance against these workflows,
please use the [conformance testing tool](https://github.com/opencontainers/distribution-spec/tree/master/conformance).

### Official Certification

Registry providers can self-cetify by submitting conformance results to
[opencontainers/oci-conformance](https://github.com/opencontainers/oci-conformance).
</speck>*/

import (
	"testing"

	g "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

func TestConformance(t *testing.T) {
	g.Describe(suiteDescription, func() {
		/*<speck tab=2>
		### Workflow Categories
		</speck>*/
		test01Pull()
		test02Push()
		test03ContentDiscovery()
		test04ContentManagement()
	})

	RegisterFailHandler(g.Fail)
	reporters := []g.Reporter{newHTMLReporter(reportHTMLFilename), reporters.NewJUnitReporter(reportJUnitFilename)}
	g.RunSpecsWithDefaultAndCustomReporters(t, suiteDescription, reporters)
}

