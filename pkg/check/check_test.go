package check

import (
	"testing"
)

func TestCheck_Cluster(t *testing.T) {
	Check_PTP()
}

func TestCheck_Python_Env(t *testing.T) {
	Check_Python_Env()
}
