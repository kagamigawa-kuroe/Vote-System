package test

import (
	"fmt"
	"ia04-vote/comsoc"
	"strconv"
	"testing"
)

func TestBordaSWF(t *testing.T) {
	prefs := [][]comsoc.Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, _ := comsoc.BordaSWF(prefs)
	fmt.Println("BordaSWF success")
	if res[1] != 4 {
		t.Errorf("error, result for 1 should be 4, %d computed", res[1])
	}
	if res[2] != 3 {
		t.Errorf("error, result for 2 should be 3, %d computed", res[2])
	}
	if res[3] != 2 {
		t.Errorf("error, result for 3 should be 2, %d computed", res[3])
	}
}

func TestBordaSCF(t *testing.T) {
	prefs := [][]comsoc.Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, err := comsoc.BordaSCF(prefs)

	if err != nil {
		t.Error(err)
	}

	if len(res) != 1 || res[0] != 1 {
		t.Errorf("error, 1 should be the only best Alternative")
	}
}

func TestMajoritySWF(t *testing.T) {
	prefs := [][]comsoc.Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, _ := comsoc.MajoritySWF(prefs)

	if res[1] != 2 {
		t.Errorf("error, result for 1 should be 2, %d computed", res[1])
	}
	if res[2] != 0 {
		t.Errorf("error, result for 2 should be 0, %d computed", res[2])
	}
	if res[3] != 1 {
		t.Errorf("error, result for 3 should be 1, %d computed", res[3])
	}
}

func TestMajoritySCF(t *testing.T) {
	prefs := [][]comsoc.Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, err := comsoc.MajoritySCF(prefs)

	if err != nil {
		t.Error(err)
	}

	if len(res) != 1 || res[0] != 1 {
		t.Errorf(strconv.FormatInt(int64(len(res)), 10))
		t.Errorf("error, 1 should be the only best Alternative")
	}
}

func TestApprovalSWF(t *testing.T) {
	prefs := [][]comsoc.Alternative{
		{1, 2, 3},
		{1, 3, 2},
		{2, 3, 1},
	}
	thresholds := []int{2, 1, 2}

	res, _ := comsoc.ApprovalSWF(prefs, thresholds)

	if res[1] != 2 {
		t.Errorf("error, result for 1 should be 2, %d computed", res[1])
	}
	if res[2] != 2 {
		t.Errorf("error, result for 2 should be 2, %d computed", res[2])
	}
	if res[3] != 1 {
		t.Errorf("error, result for 3 should be 1, %d computed", res[3])
	}
}

func TestApprovalSCF(t *testing.T) {
	prefs := [][]comsoc.Alternative{
		{1, 3, 2},
		{1, 2, 3},
		{2, 1, 3},
	}
	thresholds := []int{2, 1, 2}

	res, err := comsoc.ApprovalSCF(prefs, thresholds)

	if err != nil {
		t.Error(err)
	}
	if len(res) != 1 || res[0] != 1 {
		t.Errorf("error, 1 should be the only best Alternative")
	}
}

func TestCondorcetWinner(t *testing.T) {
	prefs1 := [][]comsoc.Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	prefs2 := [][]comsoc.Alternative{
		{1, 2, 3},
		{2, 3, 1},
		{3, 1, 2},
	}

	res1, _ := comsoc.CondorcetWinner(prefs1)
	res2, _ := comsoc.CondorcetWinner(prefs2)

	if len(res1) == 0 || res1[0] != 1 {
		t.Errorf("error, 1 should be the only best alternative for prefs1")
	}
	if len(res2) != 0 {
		t.Errorf("no best alternative for prefs2")
	}
}

func TestTieBreak1(t *testing.T) {
	a := []comsoc.Alternative{0,1,2,3,4}
	departage := comsoc.TieBreakFactory(a)
	candidat := []comsoc.Alternative{0,1}
	gagant,_ := departage(candidat)
	if gagant != 0 {
		t.Errorf("error, winner should be 0")
	}
}

func TestTieBreak2(t *testing.T) {
	a := []comsoc.Alternative{5,2,3,1,0}
	departage := comsoc.TieBreakFactory(a)

	candidat1 := []comsoc.Alternative{0,1}
	gagant1,_ := departage(candidat1)
	if gagant1 != 1 {
		t.Errorf("error, winner1 should be 0")
	}

	candidat2 := []comsoc.Alternative{3,5}
	gagant2,_ := departage(candidat2)
	if gagant2 != 5 {
		t.Errorf("error, winner should be 5")
	}
}

func TestSWFFactory(t *testing.T){
	a := []comsoc.Alternative{2,1,3}
	departage := comsoc.TieBreakFactory(a)

	prefs := [][]comsoc.Alternative{
		{1, 3, 2},
		{3, 2, 1},
		{2, 1, 3},
	}

	f := comsoc.SWFFactory(comsoc.BordaSWF,departage)
	c,_ := f(prefs)

	if c[0] != a[0] || c[1] != a[1] || c[2] != a[2]{
		t.Errorf("%d %d %d",c[0],c[1],c[2])
	}
}

func TestSCFFactory(t *testing.T){
	a := []comsoc.Alternative{1,2,3}
	departage := comsoc.TieBreakFactory(a)

	prefs := [][]comsoc.Alternative{
		{3, 1, 2},
		{3, 2, 1},
		{2, 1, 3},
	}

	f := comsoc.SCFFactory(comsoc.BordaSCF,departage)
	c,_ := f(prefs)

	if c != 3{
		t.Errorf("error of winner %d", c)
	}
}

func TestCopelandSWF(t *testing.T) {
	prefs := [][]comsoc.Alternative{
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{2, 3, 4, 1},
		{2, 3, 4, 1},
		{2, 3, 4, 1},
		{2, 3, 4, 1},
		{4, 3, 1, 2},
		{4, 3, 1, 2},
		{4, 3, 1, 2},
	}

	res, _ := comsoc.CopelandSWF(prefs)
	fmt.Println("CopelandSWF success")
	if res[1] != -1 {
		t.Errorf("error, result for -1 should be 1, %d computed", res[1])
	}
	if res[2] != 1 {
		t.Errorf("error, result for 1 should be 1, %d computed", res[2])
	}
	if res[3] != 1 {
		t.Errorf("error, result for 1 should be 1, %d computed", res[3])
	}
	if res[4] != -1 {
		t.Errorf("error, result for -1 should be 1, %d computed", res[3])
	}
}

func TestCopelandSCF(t *testing.T) {
	prefs := [][]comsoc.Alternative{
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{2, 3, 4, 1},
		{2, 3, 4, 1},
		{2, 3, 4, 1},
		{2, 3, 4, 1},
		{4, 3, 1, 2},
		{4, 3, 1, 2},
		{4, 3, 1, 2},
	}

	res, err := comsoc.CopelandSCF(prefs)

	if err != nil {
		t.Error(err)
	}

	if len(res) != 2 || res[0] != 2 || res[1] != 3 {
		t.Errorf("error, 1 should be the only best Alternative")
	}
}

func TestSTVSWF(t *testing.T) {
	prefs := [][]comsoc.Alternative{
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{2, 3, 4, 1},
		{2, 3, 4, 1},
		{2, 3, 4, 1},
		{2, 3, 4, 1},
		{4, 3, 1, 2},
		{4, 3, 1, 2},
		{4, 3, 1, 2},
	}

	res, _ := comsoc.STV_SWF(prefs)
	fmt.Println("CopelandSWF success")
	if res[1] != 1 {
		t.Errorf("error, result for -1 should be 1, %d computed", res[1])
	}
	if res[2] != -1 {
		t.Errorf("error, result for 1 should be 1, %d computed", res[2])
	}
	if res[3] != -1 {
		t.Errorf("error, result for 1 should be 1, %d computed", res[3])
	}
	if res[4] != -1 {
		t.Errorf("error, result for -1 should be 1, %d computed", res[4])
	}
}

func TestSTVSCF(t *testing.T) {
	prefs := [][]comsoc.Alternative{
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 4},
		{2, 3, 4, 1},
		{2, 3, 4, 1},
		{2, 3, 4, 1},
		{2, 3, 4, 1},
		{4, 3, 1, 2},
		{4, 3, 1, 2},
		{4, 3, 1, 2},
	}

	res, _ := comsoc.STV_SCF(prefs)
	fmt.Println("CopelandSWF success")
	if res[0] != 1 {
		t.Errorf("error, winner should be 1, %d computed", 1)
	}
}

func TestKemeny(t *testing.T) {
	prefs1 := [][]comsoc.Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}
	a,_ := comsoc.Kemeny(prefs1)
	if a[0] != 1 {
		t.Errorf("error, result for 1 should be 1, %d computed", a[1])
	}

	if a[1] != 2 {
		t.Errorf("error, result for 2 should be 2, %d computed", a[2])
	}

	if a[2] != 3 {
		t.Errorf("error, result for 3 should be 3, %d computed", a[3])
	}
}