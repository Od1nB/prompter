package git

import "testing"

func TestParseStatus(t *testing.T) {
	res := ConvStatus('A')
	if res != Added {
		t.Fatalf("got %v but expected Added", res)
	}
}
