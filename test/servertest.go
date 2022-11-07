package test

import (
	"ia04-vote/agt/sponsoragent"
	"ia04-vote/agt/voteragent"
	"ia04-vote/comsoc"
	"sync"
	"time"
)

func Test_newballot(){
	time.Sleep(2*time.Second)
	s1 := []string{"ag_id1","ag_id2","ag_id3"}
	g := sponsoragent.Sponsorinfo{"majority","Mon Jan 15:04:05 UTC 2006",s1,4}
	p := sponsoragent.Sponsoragent{g,"127.0.0.1:8082","none"}
	p.New_ballot()

	time.Sleep(2*time.Second)
	s2 := []string{"ag_id1","ag_id2","ag_id3"}
	g1 := sponsoragent.Sponsorinfo{"borda","Mon Jan 15:04:05 UTC 2006",s2,4}
	p1 := sponsoragent.Sponsoragent{g1,"127.0.0.1:8082","none"}
	p1.New_ballot()
	// a,_ := json.Marshal(g1)
	// fmt.Println(*(*string)(unsafe.Pointer(&a)))
}

func Test_vote(){
	time.Sleep(4*time.Second)
	s1 := []comsoc.Alternative{1,2,4,3}
	v1 := voteragent.Voterinfo{"ag_id1","vote1",s1,nil}
	var s sync.Mutex
	p1 := voteragent.Voteragent{s,"127.0.0.1:8082",v1}
	p1.Vote()
}

func Test_vote2(){
	time.Sleep(4*time.Second)
	s1 := []comsoc.Alternative{1,2,4,3}
	v1 := voteragent.Voterinfo{"ag_id2","vote1",s1,nil}
	var s sync.Mutex
	p1 := voteragent.Voteragent{s,"127.0.0.1:8082",v1}
	p1.Vote()
}

func Test_vote3(){
	time.Sleep(4*time.Second)
	s1 := []comsoc.Alternative{1,2,4,3}
	v1 := voteragent.Voterinfo{"ag_id3","vote1",s1,nil}
	var s sync.Mutex
	p1 := voteragent.Voteragent{s,"127.0.0.1:8082",v1}
	p1.Vote()
}

func Test_vote4(){
	time.Sleep(4*time.Second)
	s1 := []comsoc.Alternative{1,2,4,3}
	v1 := voteragent.Voterinfo{"ag_id4","vote1",s1,nil}
	var s sync.Mutex
	p1 := voteragent.Voteragent{s,"127.0.0.1:8082",v1}
	p1.Vote()

	time.Sleep(4*time.Second)
	p1.Result()
}