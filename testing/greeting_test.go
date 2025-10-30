package testmamma

import "testing"

func testcars(t *testing.T) {

	got := carboys(1000000)
	want := 700000

	if got != want {
		t.Fatalf("test case failed %q %q", got, want)
	}
}
