package tests

import "github.com/thalissonfelipe/banking/pkg/domain/vos"

// Valid CPF's used only for test purposes
var (
	TestCPF1, _ = vos.NewCPF("648.446.967-93")
	TestCPF2, _ = vos.NewCPF("626.413.228-46")
	TestCPF3, _ = vos.NewCPF("871.957.260-37")
)
