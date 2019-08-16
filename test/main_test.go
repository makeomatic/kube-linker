// TODO:
// https://github.com/operator-framework/operator-sdk/blob/master/doc/user/unit-testing.md
// https://github.com/operator-framework/operator-sdk/blob/master/doc/test-framework/writing-e2e-tests.md
package e2e

import (
	"testing"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
)

func TestMain(m *testing.M) {
	framework.MainEntry(m)
}
