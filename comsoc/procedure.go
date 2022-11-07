package comsoc

import (
	"errors"
	"sort"
)

//
// Algorithme de vote(SWF/SCF)
//

/**
 * CondorcetWinner
 * @Description: calculater gagnant de Condorcet
 * @param p: un paramètre type Profile
 * @return bestAlts: slice du gagnant
 * @return err: erreurs possibles
 */
func CondorcetWinner(p Profile) (bestAlts []Alternative, err error) {
	checkProfile(p)
	// Enregistrer le nombre de fois que chaque personne gagne pendant la comparaison par paires
	count := make(map[Alternative]int)
	note := make([]Alternative, 0)
	for i := range p[0] {
		note = append(note, p[0][i])
	}

	for i := 0; i < len(note); i++ {
		for j := i + 1; j < len(note); j++ {
			// Comparaison par paires des candidats
			a := 0
			b := 0
			for k := range p {
				index_1 := -1
				index_2 := -1
				for t := range p[k] {
					if p[k][t] == note[i] {
						index_1 = t
					}
					if p[k][t] == note[j] {
						index_2 = t
					}
				}
				if index_1 < index_2 {
					a++
				} else if index_1 > index_2 {
					b++
				}
			}
			if a > b {
				count[note[i]]++
			} else if b > a {
				count[note[j]]++
			}
		}
	}

	ans := make([]Alternative, 0)
	max_v := len(p[0]) - 1

	for i, j := range count {
		if j == max_v {
			ans = append(ans, i)
		}
	}

	if len(ans) > 1 {
		ans2 := make([]Alternative, 0)
		return ans2, nil
	} else {
		return ans, nil
	}
}

/**
 * MajoritySWF
 * @Description: SWF en Scrutin majoritaire simple
 * @param p: un paramètre type Profile
 * @return count: un paramètre type Count
 * @return err: erreurs possibles
 */
func MajoritySWF(p Profile) (count Count, err error) {
	err = checkProfile(p)
	if err != nil {
		return nil, err
	}

	count = make(Count)
	for i:=0; i<len(p[0]); i++ {
		count[p[0][i]] = 0
	}

	for i := 0; i < len(p); i++ {
		count[p[i][0]] = count[p[i][0]] + 1
	}

	return count, err
}

/**
 * MajoritySCF
 * @Description: SCF en Scrutin majoritaire simple
 * @param p: un paramètre type Profile
 * @return bestAlts: slice du gagnant
 * @return err: erreurs possibles
 */
func MajoritySCF(p Profile) (bestAlts []Alternative, err error) {
	count, e := MajoritySWF(p)
	if e == nil {
		bestAlts = maxCount(count)
	} else {
		return nil, e
	}
	return bestAlts, nil
}

/**
 * BordaSWF
 * @Description: SWF à la règle de Borda
 * @param p: un paramètre type Profile
 * @return count: un paramètre type Count
 * @return err: erreurs possibles
 */
func BordaSWF(p Profile) (count Count, err error) {
	err = checkProfile(p)
	if err != nil {
		return nil, err
	}
	count = make(Count)

	// Calculer les scores
	for i := range p {
		vote := len(p[0]) - 1
		for j := range p[i] {
			count[p[i][j]] = count[p[i][j]] + vote
			vote--
		}
	}
	return count, err
}

/**
 * BordaSCF
 * @Description: SCF à la règle de Borda
 * @param p: un paramètre type Profile
 * @return bestAlts: slice du gagnant
 * @return err: erreurs possibles
 */
func BordaSCF(p Profile) (bestAlts []Alternative, err error) {
	count, e := BordaSWF(p)
	if e == nil {
		bestAlts = maxCount(count)
		return bestAlts, nil
	} else {
		return nil, e
	}
}

/**
 * ApprovalSWF
 * @Description: SWF de Approval voting
 * @param p: un paramètre type Profile
 * @param thresholds: un slice de type int
 * @return count: un paramètre type Count
 * @return err: erreurs possibles
 */
func ApprovalSWF(p Profile, thresholds []int) (count Count, err error) {
	count = make(Count)
	for i := range p {
		for j := 0; j < thresholds[i]; j++ {
			count[p[i][j]]++
		}
	}
	return count, err
}

/**
 * ApprovalSCF
 * @Description: SCF Approval voting
 * @param p: un paramètre type Profile
 * @param thresholds: un slice de type int
 * @return bestAlts: slice du gagnant
 * @return err: erreurs possibles
 */
func ApprovalSCF(p Profile, thresholds []int) (bestAlts []Alternative, err error) {
	count, e := ApprovalSWF(p, thresholds)
	if e == nil {
		bestAlts = maxCount(count)
	} else {
		return nil, e
	}
	return bestAlts, nil
}

/**
 * KramerSimpsonSWF
 * @Description: SWF règle de KramerSimpson
 * @param p: un paramètre type Profile
 * @return Count: un paramètre type Count
 * @return error: erreurs possibles
 */
func KramerSimpsonSWF(p Profile) (Count, error) {
	note := make([]Alternative, 0)
	count := make(map[Alternative]int)
	for i := range p[0] {
		note = append(note, p[0][i])
	}
	for key, _ := range note {
		count[note[key]] = 0
	}
	for i := 0; i < len(note); i++ { //Comparer les voix du candidat i et du candidat j
		for j := i + 1; j < len(note); j++ {
			a := 0
			b := 0
			for k := range p {
				index_1 := -1
				index_2 := -1
				for t := range p[k] {
					if p[k][t] == note[i] {
						index_1 = t
					}
					if p[k][t] == note[j] {
						index_2 = t
					}
				}
				if index_1 < index_2 {
					a++ // On augmente la valeur de a par 1 si le candidat i est préféré que le candidat j
				}
				if index_1 > index_2 {
					b++ // On augmente la valeur de b par 1 si le candidat j est préféré que le candidat i
				}
			}
			//Mise à jour de la note de chaque candidat : on choisit le minimum
			if count[note[i]] == 0 {
				count[note[i]] = a
			} else if count[note[i]] > a {
				count[note[i]] = a
			}

			if count[note[j]] == 0 {
				count[note[j]] = b
			} else if count[note[j]] > b {
				count[note[j]] = b
			}

		}

	}
	return count, nil
}

/**
 * KramerSimpsonSCF
 * @Description: SCF KramerSimpson
 * @param p: un paramètre type Profile
 * @return bestAlts: slice du gagnant
 * @return err: erreurs possibles
 */
func KramerSimpsonSCF(p Profile) (bestAlts []Alternative, err error) {
	max_index := 0
	count, err := KramerSimpsonSWF(p)
	if err != nil {
		return nil, err
	}
	for key, value := range count {
		if value > count[Alternative(max_index)] {
			max_index = int(key)
		}
	}
	bestAlts = append(bestAlts, Alternative(max_index))
	return bestAlts, nil
}


/**
 * CopelandSWF
 * @Description: SCF règle de Copeland
 * @param p: un paramètre type Profile
 * @return Count: un paramètre type Count
 * @return error: erreurs possibles
 */
func CopelandSWF(p Profile) (Count, error) {
	count := make(map[Alternative]int)
	note := make([]Alternative, 0)
	for i := range p[0] {
		note = append(note, p[0][i])
	}

	for i := 0; i < len(note); i++ {
		for j := i + 1; j < len(note); j++ {
			a := 0
			b := 0
			for k := range p {
				index_1 := -1
				index_2 := -1
				for t := range p[k] {
					if p[k][t] == note[i] {
						index_1 = t
					}
					if p[k][t] == note[j] {
						index_2 = t
					}
				}
				if index_1 < index_2 {
					a++
				} else if index_1 > index_2 {
					b++
				}
			}
			if a > b {
				count[note[i]]++
				// Une soustraction supplémentaire de plus que Condorcet
				count[note[j]]--
			} else if b > a {
				count[note[j]]++
				count[note[i]]--
			}
		}
	}
	return count, nil
}

/**
 * CopelandSCF
 * @Description: SCF Copeland
 * @param p: un paramètre type Profile
 * @return bestAlts: slice du gagnant
 * @return err: erreurs possibles
 */
func CopelandSCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := CopelandSWF(p)
	if err != nil {
		return nil, err
	}
	ans := make([]Alternative, 0)
	max_v := -1
	for _, j := range count {
		if j > max_v {
			max_v = j
		}
	}

	for i, j := range count {
		if j == max_v {
			ans = append(ans, i)
		}
	}
	return ans, nil
}

/**
 * Coombs_SWF
 * @Description: STV Coombs SWF
 * @param p: un paramètre type Profile
 * @return count: score des candidats, 1 signifie gagnant, 0 signifie d'etre éliminé
 * @return err: erreurs possibles
 */
func CoombsSWF(p Profile) (count Count, err error) {
	count = make(Count)
	err = checkProfile(p)
	if err != nil {
		return nil, err
	}

	// supprimer un candidat
	del := func(data *Profile, a Alternative) {
		for i := 0; i < len(*data); i++ {
			j := 0
			for ; (*data)[i][j] != a; j++ {
			}
			k := j + 1
			for ; k < len((*data)[i]); k++ {
				(*data)[i][k-1] = (*data)[i][k]
			}
			(*data)[i][k-1] = -1
		}
	}
	size := len(p[0])
	for i := range p[0] {
		count[p[0][i]] = 1
	}

	for i := 0; i < size-1; i++ {
		max_candidat := p[0][0]
		m := make(map[Alternative]int)
		for j := 0; j < len(p); j++ {
			m[p[j][size-1-i]]++
		}
		for key, value := range m {
			if value > m[max_candidat] {
				max_candidat = key
			}

		}
		count[max_candidat] = -1
		del(&p, max_candidat)
	}

	return count, nil
}


/**
 * CoombsSCF
 * @Description: STV Coombs SCF
 * @param p: un paramètre type Profile
 * @return bestAlts: slice de gagnants
 * @return err: erreurs possibles
 */
func CoombsSCF(p Profile) (bestAlts []Alternative, err error) {
	c, err := CoombsSWF(p)
	if err != nil {
		return nil, err
	}
	for i, j := range c {
		if j == 1 {
			bestAlts = append(bestAlts, Alternative(i))
		}
	}
	return bestAlts, nil
}

/**
 * STV_SWF
 * @Description: Vote Simple Transférable SWF
 * @param p: un paramètre type Profile
 * @return count: score des candidats, 1 signifie gagnant, 0 signifie d'etre éliminé
 * @return err: erreurs possibles
 */
func STV_SWF(p Profile) (count Count, err error) {
	err = checkProfile(p)
	if err != nil {
		return nil, err
	}
	count = make(Count)

	// supprimer un candidat
	del := func(data *Profile, a Alternative) {
		for i := 0; i < len(*data); i++ {
			j := 0
			for ; (*data)[i][j] != a; j++ {
			}
			k := j + 1
			for ; k < len((*data)[i]); k++ {
				(*data)[i][k-1] = (*data)[i][k]
			}
			(*data)[i][k-1] = -1
		}
	}

	size := len(p[0])
	for i := range p[0] {
		count[p[0][i]] = 1
	}

	for i := 0; i < size-1; i++ {
		var a Alternative = p[0][0]
		m := make(map[Alternative]int)
		for j := 0; j < len(p); j++ {
			m[p[j][0]]++
		}
		for k, j := range m {
			if j < m[a] {
				a = k
			}
		}
		count[a] = -1
		del(&p, a)
	}

	return count, err
}

/**
 * STV_SCF
 * @Description: Vote Simple Transférable SCF
 * @param p: un paramètre type Profile
 * @return bestAlts: slice de gagnants
 * @return err: erreurs possibles
 */
func STV_SCF(p Profile) (bestAlts []Alternative, err error) {
	c, err := STV_SWF(p)
	if err != nil {
		return nil, err
	}
	for i, j := range c {
		if j == 1 {
			bestAlts = append(bestAlts, i)
		}
	}
	return bestAlts, nil
}

/**
 * Kemeny_SWF
 * @Description: Kemeny vote SWF
 * @param p: un paramètre type Profile
 * @return ans: gagnant
 * @return e: erreurs possibles
 */
func Kemeny_SWF(p Profile) ( ans []Alternative, e error) {
	e = checkProfile(p)
	if e != nil {
		return nil, e
	}
	// tous alternative
	all := make([]Alternative, 0)
	for i := range p[0] {
		all = append(all, p[0][i])
	}

	ans = nil
	min_distance := int(^uint(0) >> 1)
	all_combination := Permute(all)

	for _, j := range all_combination {
		d, err := Distance_edit_somme(j, p)
		if err != nil {
			e = err
			return nil, e
		}
		if d < min_distance {
			ans = j
			min_distance = d
		}
	}

	return ans, nil
}

/**
 * Kemeny
 * @Description: Kemeny vote SCF
 * @param p: un paramètre type Profile
 * @return ans: ordre de candidats
 * @return e: erreurs possibles
 */
func Kemeny_SCF(p Profile) (ans []Alternative, e error) {
	a,e := Kemeny_SWF(p)
	if e != nil {
		return nil, e
	}
	temp := []Alternative{a[0]}
	return temp,nil
}

func SinglePeakSWF(p Profile) (count Count, err error) {
	count = make(Count)
	//for i := range p { // Vérification de single-peaked
	//	order_inf := make([]Alternative, 0) // Stocker les candidats de l'ordre inférieur que le premier
	//	order_sup := make([]Alternative, 0) // Stocker les candidats de l'ordre supérieur que le premier
	//	med := p[i][0]
	//	order_sup = append(order_sup, med)
	//	order_inf = append(order_inf, med)
	//	for j := 1; j < len(p[i]); j++ {
	//		if p[i][j] < med { // Stocker les candidats de l'ordre inférieur que le premier dans ordre_inf
	//			order_inf = append(order_inf, p[i][j])
	//		} else { // Stocker les candidats de l'ordre supérieur que le premier dans ordre_sup
	//			order_sup = append(order_sup, p[i][j])
	//		}
	//	}
	//	for k := 0; k < len(order_sup)-1; k++ { // Vérification de l'ordre stricte
	//		if order_sup[k] > order_sup[k+1] {
	//			return nil, errors.New("Pas Single-Peaked")
	//		}
	//	}
	//	for k := 0; k < len(order_inf)-1; k++ {
	//		if order_inf[k] < order_inf[k+1] {
	//			return nil, errors.New("Pas Single-Peaked")
	//		}
	//	}
	//}

	for _, value := range p[0] {
		count[value] = 0
	}
	for i := range p {
		count[p[i][0]]++
	}
	return count, nil
}

func SinglePeakedSCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := SinglePeakSWF(p)
	note := make([]int, 0)
	var median int
	for _, elem := range count {
		note = append(note, elem)
	}
	length := len(note)
	sort.Ints(note)
	if length%2 == 0 {
		median = note[int(length/2-1)]
	} else {
		median = note[int(length/2)]
	}
	for key, elem := range count {
		if elem == median {
			bestAlts = append(bestAlts, key)
		}
	}
	return bestAlts, nil
}


//
// Factory function
//

/**
 * TieBreakFactory
 * @Description: générer une fonction de comparaison de départage
 * @param a: priorité des candidats au départage
 * @return tiebreak: fonction de départage générée
 */
func TieBreakFactory(a []Alternative) (tiebreak func([]Alternative) (Alternative, error)) {
	note := make(map[Alternative]int)
	for i := range a {
		note[a[i]] = i
	}
	return func(a []Alternative) (alt Alternative, e error) {
		if len(a) == 0 {
			return -1, errors.New("alternative list est Null")
		}
		var ans Alternative = a[0]
		for i := range a {
			if note[a[i]] < note[ans] {
				ans = a[i]
			}
		}
		return ans, nil
	}
}

/**
 * SWFFactory
 * @Description: combiner les fonctions de vote SWF avec TieBreak
 * @param s: la fonction de vote
 * @param t: priorité des candidats au départage
 * @return swf: priorité de tous les candidat en départage
 */
func SWFFactory(s func(p Profile) (Count, error), t func([]Alternative) (Alternative, error)) (swf func(Profile) ([]Alternative, error)) {
	return func(p Profile) ([]Alternative, error) {
		temp, e := s(p)
		if e != nil {
			return nil, e
		}

		type candi struct {
			a     Alternative
			score int
		}

		note := make([]candi, 0)

		for i, j := range temp {
			n := candi{i, j}
			note = append(note, n)
		}

		sort.Slice(note, func(i, j int) bool {
			l := []Alternative{note[i].a, note[j].a}
			pre, _ := t(l)
			return note[i].score > note[j].score || (note[i].score == note[j].score && pre == note[i].a)
		})

		ans := make([]Alternative, 0)
		for _, j := range note {
			ans = append(ans, j.a)
		}

		return ans, nil
	}
}

/**
 * SCFFactory
 * @Description: combiner les fonctions de vote SCF avec TieBreak
 * @param s: la fonction de vote
 * @param t: priorité des candidats au départage
 * @return scf: vainqueur du départage
 */
func SCFFactory(s func(p Profile) ([]Alternative, error), t func([]Alternative) (Alternative, error)) (scf func(Profile) (Alternative, error)) {
	return func(p Profile) (Alternative, error) {
		temp, e := s(p)
		if e != nil {
			return -1, e
		}
		a, err := t(temp)
		if err != nil {
			return -1, err
		}
		return a, nil
	}
}
