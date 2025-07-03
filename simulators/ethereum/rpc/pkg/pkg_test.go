package pkg_test

import (
	"testing"

	"github.com/ethereum/hive/simulators/ethereum/rpc/pkg"
)

func TestRunAl(t *testing.T) {
	pkg.RunAllTests(t, "localhost", "wasp")
}
