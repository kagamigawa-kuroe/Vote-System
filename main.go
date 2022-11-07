package main

import "gitlab.utc.fr/wanhongz/ia04-vote/agt/ballotagent"

/**
 * main
 * @Description: Démarrer le serveur de vote
 */
func main() {
	ballotagent.StartVoteServer()
}