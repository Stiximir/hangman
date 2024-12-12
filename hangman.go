package hangman

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

type Data struct { // Structure de données pour le jeu
	Life        int
	Word        string
	HiddenWord  string
	Letter      string
	AlreadyUsed []string
}

func load(savefile string) *Data {
	var jeu Data
	game := readWordsFromFile(savefile) // lit le fichier de sauvegarde
	if len(game) == 0 {
		log.Fatal("Empty File.")
	}
	err := json.Unmarshal([]byte(game[0]), &jeu) // récupère les valeurs de jeu
	if err != nil {
		log.Fatal(err)
	}
	return &jeu
}

func save(jeu *Data) {
	game, err1 := json.Marshal(jeu) // convertit Game en JSON
	if err1 != nil {
		log.Fatal(err1)
	}
	err2 := os.WriteFile("Data/game.txt", game, 0644) // crée ou écrase partie.txt avec partie
	if err2 != nil {
		log.Fatal(err2)
	}
}

func readWordsFromFile(filename string) []string { // Lecture du fichier, retourne la liste de mots avec words
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return words
}

func HideWordStart(jeu *Data, randWord string) {
	var doubleCheck bool
	jeu.Word = randWord // Le mot complet dans jeu.Word
	hiddenword := ""    // Le mot incomplet qui va être construit
	rand.Seed(time.Now().UnixNano())
	var randomletter []int
	for i := 0; i < len(randWord)/2-1; i++ { // On recherche des index aléatoires du mot complet
		randomletter = append(randomletter, rand.Intn(len(randWord)))
		if len(randomletter) > 1 {
			for n := 0; n < len(randomletter); n++ {
				for randWord[randomletter[i]] == randWord[randomletter[n]] && i != n {
					randomletter[i] = rand.Intn(len(randWord))
				}
			}
		}
	}
	var randomlettervalue []rune
	for i, v := range randWord { // On convertit ces index en des lettres pour les afficher plus facilement
		for n := 0; n < len(randomletter); n++ {
			if randomletter[n] == i {
				randomlettervalue = append(randomlettervalue, v)
			}
		}
	}
	for _, v := range randWord {
		doubleCheck = true
		reveal := false
		for _, n := range randomlettervalue {
			if n == v {
				reveal = true
				break
			}
		}
		if reveal {
			hiddenword = hiddenword + string(v)
			if len(jeu.AlreadyUsed) == 0 {
				jeu.AlreadyUsed = append(jeu.AlreadyUsed, string(v))
			} else {
				for _, n := range jeu.AlreadyUsed {
					if n == string(v) {
						doubleCheck = false
					}
				}
			}
			if doubleCheck {
				jeu.AlreadyUsed = append(jeu.AlreadyUsed, string(v))
			}
		} else {
			hiddenword = hiddenword + "_"
		}
	}
	jeu.HiddenWord = hiddenword // Intégration du mot caché dans les données
}

func DisplayHangman(jeu *Data) string {
	hangman := readWordsFromFile("Data/hangman.txt")
	var linenumber int = 8
	linenumber = linenumber * (10 - jeu.Life - 1)
	var hangmanDrawing string
	for i := 0; i < 7; i++ {
		if linenumber >= 80 || jeu.Life == 10 {
			break
		}
		hangmanDrawing += hangman[linenumber][0:9] + "\n"
		linenumber++
	}
	return hangmanDrawing
}

func lives(jeu *Data) {
	jeu.Life -= 1
	if len(jeu.Letter) > 1 && jeu.Letter != jeu.Word && jeu.Life != 1 {
		jeu.Life -= 1
	}
}

func VictoryCheck(jeu *Data) bool {
	return jeu.Word == jeu.HiddenWord || jeu.Letter == jeu.Word
}

func UpdWord(jeu *Data) string {
	hiddenRunes := []rune(jeu.HiddenWord) // Ajout de chacune de ces valeurs en rune pour des raisons d'écriture de string plus simple
	wordRunes := []rune(jeu.Word)
	letterRune := rune(jeu.Letter[0])
	var success bool = false
	for i := 0; i < len(jeu.Word); i++ { // Vérification pour une lettre si elle est présente dans le mot
		if wordRunes[i] == letterRune && len(jeu.Letter) == 1 {
			hiddenRunes[i] = wordRunes[i]
			success = true
		}
	}
	if success {
	}
	if !success || (len(jeu.Letter) > 1 && jeu.Letter != jeu.Word) {
		lives(jeu)
	}
	if VictoryCheck(jeu) {
		return fmt.Sprintf("Congratulations, you found %s! You had %v lives remaining.", jeu.Word, jeu.Life)
	}
	jeu.HiddenWord = string(hiddenRunes)
	return jeu.HiddenWord
}

func RandomWord(filename string) string {
	words := readWordsFromFile(filename)
	rand.Seed(time.Now().UnixNano())
	randomWord := words[rand.Intn(len(words))]
	return randomWord
}

func GameStart(level string) *Data {
	var jeu Data
	jeu.Life = 10
	var filename string
	switch level {
	case "1":
		filename = "Data/words.txt"
	case "2":
		filename = "Data/words2.txt"
	case "3":
		filename = "Data/words3.txt"
	default:
		filename = "Data/words.txt"
	}
	HideWordStart(&jeu, RandomWord(filename))
	return &jeu
}

// le hangman fonctionne uniquement avec la version WEB
