package actiontest

// Actions []name.Ident
// Name (name.Ident)
// TestPkg string
const actionsTestTmpl = `
package {{.TestPkg}}

import (
  "testing"

  "github.com/stretchr/testify/require"
)
{{ range $a := .Actions }}
func (as *ActionSuite) Test_{{$.Name.Resource}}Resource_{{ $a.Pascalize }}() {
  as.Fail("Not Implemented!")
}
{{ end }}
`
