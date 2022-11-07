package test

import (
	"ia04-vote/comsoc"
	"testing"
)

func TestPermutation(t *testing.T) {
	a := []comsoc.Alternative{1,2,3}
	per := comsoc.Permute(a)
	if len(per) != 6 {
		t.Errorf("error, result should be 6")
	}
}