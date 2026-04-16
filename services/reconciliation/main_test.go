package main

import "testing"

func TestRunReconciliation(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("RunReconciliation panicked: %v", r)
		}
	}()

	RunReconciliation()
}

func TestGenerateAuditReport(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("GenerateAuditReport panicked: %v", r)
		}
	}()

	GenerateAuditReport()
}
