package ballotagent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitlab.utc.fr/wanhongz/ia04-vote/agt/sponsoragent"
	"gitlab.utc.fr/wanhongz/ia04-vote/agt/voteragent"
	"gitlab.utc.fr/wanhongz/ia04-vote/comsoc"
	io "io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

/**
 * Ballotagentmanager
 * @Description: Gérer tous les votes, main server
 */
type Ballotagentmanager struct {
	sync.Mutex
	IP           string  // ip du serveur
	Port         string  // port du serveur
	Ballotagents map[string]*Ballotagent // map du tous les Ballotagents
	NowID        int // prochain identifiant disponible de vote
}

/**
 * handlerNewBallot
 * @Description: Le processeur de la requête /new_ballot
 */
func (bs *Ballotagentmanager) handlerNewBallot(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime)
	log.Println(": Get a new ballot create request")

	var resp sponsoragent.Response
	var re sponsoragent.Sponsorinfo

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err := json.Unmarshal(buf.Bytes(), &re)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		var a sync.Mutex
		id := "vote"

		bs.Lock()
		id += strconv.Itoa(bs.NowID)

		// Vérifier l'heure d'échéance
		t1, e1 := time.ParseInLocation("Mon Jan 2 15:04:05 UTC 2006", re.Deadline, time.Local)
		t2, _ := time.ParseInLocation("Mon Jan 2 15:04:05 UTC 2006", time.Now().Format("Mon Jan 2 15:04:05 UTC 2006"), time.Local)
		if e1 != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(": Fail to create a new ballot " + resp.ID + ", because time is not valid")
			bs.Unlock()
			return
		}

		ex := t1.Unix() - t2.Unix()
		if ex <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(": Fail to create a new ballot " + resp.ID + ", because time is not valid")
			bs.Unlock()
			return
		}

		if _, ok := method_scf[re.Rule]; ok {
			// Créer un nouveau vote
			b := &Ballotagent{a, re, make([]voteragent.Voterinfo, 0),
				make(map[string]bool), make(comsoc.Profile, 0), id, false, int(ex), make([]int, 0)}

			// timer
			go b.SetFinished()

			// préparer tous les votants
			for i := 1; i <= len(re.Voter_ids); i++ {
				id := "ag_id"
				id += strconv.Itoa(i)
				b.Voters[id] = true
			}
			bs.Ballotagents[id] = b
			resp.ID = id
			bs.NowID++

			// retourner
			w.WriteHeader(http.StatusCreated)
			log.Println(": Create a new ballot " + resp.ID)
			serial, _ := json.Marshal(resp)
			w.Write(serial)
			bs.Unlock()
		} else {
			w.WriteHeader(http.StatusNotImplemented)
			log.Println(": Fail to create a new ballot " + resp.ID + ", vote method is note valid")
			bs.Unlock()
		}
	}
}

/**
 * handlerVoteRequest
 * @Description: Le processeur de la requête /vote
 */
func (bs *Ballotagentmanager) handlerVoteRequest(w http.ResponseWriter, r *http.Request) {
	bs.Lock()
	log.SetFlags(log.Ldate | log.Ltime)

	var re voteragent.Voterinfo
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err := json.Unmarshal(buf.Bytes(), &re)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		bs.Unlock()
		return
	}

	if _, ok := bs.Ballotagents[re.Vote_ID]; ok {
		agent := bs.Ballotagents[re.Vote_ID]
		agent.getNewVoteRequest(re, w)
		bs.Unlock()
	} else {
		w.WriteHeader(http.StatusNotImplemented)
		log.Println(": Get a new vote request of " + re.Vote_ID + ", from " + re.Agent_ID + ", " +
			"vote failed because the " + re.Vote_ID + " do not exist")
		bs.Unlock()
		return
	}
}


/**
 * handlerResultRequest
 * @Description: Le processeur de la requête /result
 */
func (bs *Ballotagentmanager) handlerResultRequest(w http.ResponseWriter, r *http.Request) {
	bs.Lock()
	log.SetFlags(log.Ldate | log.Ltime)
	var resp voteragent.Response_Result
	var re voteragent.Request_Result
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err := json.Unmarshal(buf.Bytes(), &re)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		serial, _ := json.Marshal(resp)
		w.Write(serial)
		bs.Unlock()
		return
	}
	log.Println(": Get a new result request of " + re.Ballot_Id)
	var agent *Ballotagent
	if _, ok := bs.Ballotagents[re.Ballot_Id]; ok {
		agent = bs.Ballotagents[re.Ballot_Id]
		agent.getNewResultRequest(re.Ballot_Id, w)
		bs.Unlock()
	} else {
		w.WriteHeader(http.StatusNotFound)
		serial, _ := json.Marshal(resp)
		w.Write(serial)
		bs.Unlock()
		return
	}
}

/**
 * Start
 * @Description: fonction de démarrage du serveur
 */
func (bs *Ballotagentmanager) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/new_ballot", bs.handlerNewBallot)
	mux.HandleFunc("/vote", bs.handlerVoteRequest)
	mux.HandleFunc("/result", bs.handlerResultRequest)

	s := &http.Server{
		Addr:           "127.0.0.1:8082",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.SetFlags(log.Ldate | log.Ltime)
	log.Println(": Start listen on \"" + bs.IP + ":" + bs.Port + "\"")
	log.Fatal(s.ListenAndServe())
}

type ServerInfo struct {
	IP  string `json:ip`
	Port  string `json:port`
}

var file_locker sync.Mutex

func LoadConfig(filename string) (ServerInfo, bool) {
	log.SetFlags(log.Ldate | log.Ltime)
	var conf ServerInfo
	file_locker.Lock()
	data, err := io.ReadFile(filename) //read config file
	file_locker.Unlock()
	if err != nil {
		log.Println("read json file error")
		return conf, false
	}
	datajson := []byte(data)
	err = json.Unmarshal(datajson, &conf)
	if err != nil {
		log.Println("unmarshal json file error")
		return conf, false
	}
	return conf, true
}

func InitConfig() ServerInfo {
	log.SetFlags(log.Ldate | log.Ltime)
	conf, bl := LoadConfig("./config.json") //get config struct
	if !bl {
		log.Println("Init config file failed")
		return ServerInfo{"127.0.0.1","8082"}
	}
	log.Println(": Init config file successed")
	return conf
}

/**
 * StartVoteServer
 * @Description: L'interface pour démarrer la fonction
 * @param IP：ip de serveur
 * @param Port: port de serveur
 */
func StartVoteServer() {
	banner := "  ___    _    ___  _  _      __     __    _       \n " +
		"|_ _|  / \\  / _ \\| || |     \\ \\   / /__ | |_ ___ \n  " +
		"| |  / _ \\| | | | || |_ ____\\ \\ / / _ \\| __/ _ \\\n  " +
		"| | / ___ \\ |_| |__   _|_____\\ V / (_) | ||  __/\n " +
		"|___/_/   \\_\\___/   |_|        \\_/ \\___/ \\__\\___| \n"
	fmt.Println(banner)
	log.SetFlags(log.Ldate | log.Ltime)
	var mutex sync.Mutex
	s := InitConfig()
	bs := Ballotagentmanager{mutex, s.IP, s.Port, make(map[string]*Ballotagent), 0}
	bs.Start()
}