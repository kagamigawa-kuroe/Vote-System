package sponsoragent

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

/**
 * Sponsorinfo
 * @Description: Paramètres requis pour la demande /new_ballot
 */
type Sponsorinfo struct {
	Rule string `json:"rule"` // "majority","borda", "approval", "stv", "kemeny",...
	Deadline string `json:"deadline"`// e.g. "Tue Nov 10 23:00:00 UTC 2009"
	Voter_ids []string  `json:"voter_Ids"`// e.g. ["ag_id1", "ag_id2", "ag_id3"]
	Alts int `json:"alts"`
}

/**
 * Sponsoragent
 * @Description: agent pour initier un vote
 */
type Sponsoragent struct {
	Sponsorinfo
	ServerAddress string
	ID string //L'identifiant doit être défini "none" lors de l'initialisation de l'agent.
}

// tous les vote algorithme
var method = map[string]int {
	"condorcet"     :0,
	"majority"      :1,
	"borda"         :2,
	"kramersimpson" :3,
	"copeland"      :4,
	"coombs"        :5,
	"stv"           :6,
	"kemeny"        :7,
	"singlepeak"    :8,
}

type Response struct {
	ID string `json:"id"`
}

/**
 * New_ballot
 * @Description: envoyer une demande de /new_ballot
 * @return error: erreurs possibles
 */
func (s *Sponsoragent) New_ballot() error{
	if s.ID != "none" {
		return errors.New("this agent has already been registered")
	}

	req := Sponsorinfo{
		s.Rule,s.Deadline,s.Voter_ids,s.Alts,
	}

	url := "http://" + s.ServerAddress + "/new_ballot"

	data, e := json.Marshal(req)
	if e != nil {
		return e
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	var re Response
	json.Unmarshal(buf.Bytes(), &re)
	if resp.StatusCode == http.StatusCreated {
		log.Println(": a new vote "+ re.ID + " create successfully")
		s.ID = re.ID
		return nil
	}else if resp.StatusCode == http.StatusBadRequest {
		log.Println("server had a bad request, registration failed")
		return errors.New("server had a bad request, registration failed")
	}else {
		log.Println("server had not implemented this function")
		return errors.New("server had not implemented this function")
	}
	return nil
}