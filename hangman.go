package hangman

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type Data struct { // Structure de données pour le jeu
	Life        int
	Word        string
	HiddenWord  string
	Letter      string
	AlreadyUsed []string
}

func consoleReset() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func load(savefile string) {
	var jeu Data
	game := readWordsFromFile(savefile) // lit le fichier de sauvegarde
	if len(game) == 0 {
		log.Fatal("Empty File.")
	}
	err := json.Unmarshal([]byte(game[0]), &jeu) // récupère les valeurs de jeu
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Welcome back, you have %d tries remaining.\n\n", jeu.Life)
	for i := 0; i < len(jeu.HiddenWord); i++ {
		fmt.Print(string(jeu.HiddenWord[i]), " ")
	}
	fmt.Println()
	input(&jeu)
}

func save(jeu *Data) {
	game, err1 := json.Marshal(jeu) // convertit Game en JSON
	if err1 != nil {
		log.Fatal(err1)
	}
	err2 := os.WriteFile("game.txt", game, 0644) // crée ou écrase partie.txt avec partie
	if err2 != nil {
		log.Fatal(err2)
	}
}

func readWordsFromFile(filename string) []string { // Lecture du fichier, retourne la liste de mot avec words
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

func hideWordStart(jeu *Data, randWord string) { // Mise en place des variables Word et Hiddenword de la structure jeu
	jeu.Word = randWord // Le mot complet dans jeu.Word
	hiddenword := ""    // le mot incomplet qui va etre construit
	rand.Seed(time.Now().UnixNano())
	var randomletter []int
	for i := 0; i < len(randWord)/2-1; i++ { // On recherche des index aleatoire du mot complet selon la taille du mot
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
	for i, v := range randWord { // On convertis ces index en des lettres pour pouvoir toutes les scanner et les placer dans la prochaine boucle
		for n := 0; n < len(randomletter); n++ {
			if randomletter[n] == i {
				randomlettervalue = append(randomlettervalue, v)
			}
		}
	}
	for i, v := range randWord { // Vérification a chaque tour que les lettres de la liste corresponde a la lettre que l'on observe actuellement dans la boucle
		reveal := false
		for _, n := range randomlettervalue {
			if n == v {
				reveal = true
				break
			}
		}
		if reveal { // Si lettre trouvé
			hiddenword = hiddenword + string(v)
			jeu.AlreadyUsed = append(jeu.AlreadyUsed, string(v))
			fmt.Print(string(hiddenword[i]), " ")
		} else { // Sinon placement du cache
			hiddenword = hiddenword + "_"
			fmt.Print(string(hiddenword[i]), " ")
		}
	}
	jeu.HiddenWord = hiddenword // Intégration du mot caché dans les données
	fmt.Println()
	input(jeu) // Requete au joueur
}

func DisplayHangman(jeu *Data) {
	hangman := readWordsFromFile("hangman.txt")
	var linenumber int = 8
	linenumber = linenumber * (10 - jeu.Life - 1)
	fmt.Println()
	for i := 0; i < 7; i++ {
		if linenumber >= 80 || jeu.Life == 10 {
			break
		}
		fmt.Println(hangman[linenumber][0:9])
		linenumber++
	}
}

func lives(jeu *Data) { // Gestion des vies et de la lose
	jeu.Life -= 1
	if len(jeu.Letter) > 1 && jeu.Letter != jeu.Word && jeu.Life != 1 {
		jeu.Life -= 1
		fmt.Printf("You lost 2 lives !\n\n")
	}
	if jeu.Life <= 0 {
		consoleReset()
		fmt.Println("You lost, the word was", jeu.Word, "\n")
	} else {
		fmt.Printf("You have %v lives remaining\n\n", jeu.Life)
		for i := 0; i < len(jeu.HiddenWord); i++ {
			fmt.Print(string(jeu.HiddenWord[i]), " ")
		}
		fmt.Println()
	}
}

func victoryCheck(jeu *Data) bool { // Gestion de la victoire
	if jeu.Word == jeu.HiddenWord || jeu.Letter == jeu.Word {
		return true
	}
	return false
}

func updWord(jeu *Data) { // Update le mot a chaque requete du joueur
	consoleReset()
	hiddenRunes := []rune(jeu.HiddenWord) // Ajout de chacune de ces valeurs en rune pour des raisons d'écriture de string plus simple
	wordRunes := []rune(jeu.Word)
	letterRune := rune(jeu.Letter[0])
	var success bool = false
	for i := 0; i < len(jeu.Word); i++ { // Vérification pour une lettre si elle est présente dans le mot, puis remplacement
		if wordRunes[i] == letterRune && len(jeu.Letter) == 1 {
			hiddenRunes[i] = wordRunes[i]
			success = true
		}
	}
	if success {
		fmt.Printf("You have found a letter !\n\n")
	}
	if !success || (len(jeu.Letter) > 1 && jeu.Letter != jeu.Word) { // Vérification si la vérif de la lettre est pas vrai ou si on vérif un mot et que le mot est faux
		lives(jeu)
	} else if !victoryCheck(jeu) { // Réecriture du mot caché
		jeu.HiddenWord = string(hiddenRunes)
		for i := 0; i < len(jeu.HiddenWord); i++ {
			fmt.Print(string(jeu.HiddenWord[i]), " ")
		}
		fmt.Println()
	}
	if victoryCheck(jeu) { // Si on a gagné
		consoleReset()
		fmt.Println("Congratulations, you found", jeu.Word, "\n")
		DisplayHangman(jeu)
		fmt.Printf("\nYou had %v lives remaining\n", jeu.Life)
	} else { // Sinon on redemande la une lettre
		input(jeu)
	}
}

func input(jeu *Data) { // Demander au joueur une lettre
	var letter string
	var savestr string
	var minimalize rune
	DisplayHangman(jeu)
	jeu.Letter = ""
	if jeu.Life > 0 {
		fmt.Println()
		fmt.Println("Type a letter: ")
		fmt.Scan(&letter)
		for i := 0; i < len(jeu.AlreadyUsed); i++ { // Vérification si lettre déja écrite ou présente
			if jeu.AlreadyUsed[i] == letter {
				consoleReset()
				for i := 0; i < len(jeu.HiddenWord); i++ {
					fmt.Print(string(jeu.HiddenWord[i]), " ")
				}
				fmt.Printf("\n\nYou've already written this letter or those letters\n\nYou have %v lives remaining\n", jeu.Life)
				input(jeu)
				return
			}
		}
		if letter == "STOP" { // Possibilité de stopper le jeu
			consoleReset()
			fmt.Printf("Game stop !\n\nDo you want to save your game ?\nFor yes, print 'y'\nFor no, print anything else\n")
			fmt.Scan(&savestr)
			if savestr == "y" {
				save(jeu)
				consoleReset()
				fmt.Printf("Game saved !\n\nTo restart, use '--startWith game.txt' option\n")
			}
			os.Exit(0)
		} else {
			jeu.AlreadyUsed = append(jeu.AlreadyUsed, letter)
		}
		if len(letter) > 1 { // Vérification pour les mots entiers
			for i := 0; i < len(letter); i++ {
				if rune(letter[i]) >= 'A' && rune(letter[i]) <= 'Z' {
					minimalize = rune(letter[i]) + 32
					jeu.Letter = jeu.Letter + string(minimalize)
					fmt.Print(jeu.Letter)
				} else if rune(letter[i]) >= 'a' && rune(letter[i]) <= 'z' {
					jeu.Letter = jeu.Letter + string(letter[i])
					fmt.Print(jeu.Letter)
				} else if rune(letter[i]) < 'a' || rune(letter[i]) > 'z' { // Vérifications que sa soit des lettres
					consoleReset()
					fmt.Printf("Error, try again\n\nYou have %v lives remaining\n\n", jeu.Life)
					for i := 0; i < len(jeu.HiddenWord); i++ {
						fmt.Print(string(jeu.HiddenWord[i]), " ")
					}
					fmt.Println()
					input(jeu)
					break
				}
			}
		} else if rune(letter[0]) >= 'A' && rune(letter[0]) <= 'Z' && len(letter) == 1 { // Vérification pour les lettres
			minimalize = rune(letter[0]) + 32
			letter = string(minimalize)
			jeu.Letter = letter
		} else if rune(letter[0]) >= 'a' && rune(letter[0]) <= 'z' && len(letter) == 1 {
			jeu.Letter = letter
		} else if rune(letter[0]) < 'a' || rune(letter[0]) > 'z' && len(letter) == 1 { // Vérifications que sa soit une lettre
			consoleReset()
			fmt.Printf("Error, try again\n\nYou have %v lives remaining\n\n", jeu.Life)
			for i := 0; i < len(jeu.HiddenWord); i++ {
				fmt.Print(string(jeu.HiddenWord[i]), " ")
			}
			fmt.Println()
			input(jeu)
		}
		updWord(jeu)
	}
}

func randomWord(filename string) string { // Prendre aléatoirement un mot de la liste choisis en paramètre
	words := readWordsFromFile(filename)
	rand.Seed(time.Now().UnixNano())
	randomWord := words[rand.Intn(len(words))]
	return randomWord
}

func GameStart() { // Démarrage du jeu
	consoleReset()
	if len(os.Args) == 3 && os.Args[1] == "--startWith" {
		load(os.Args[2])
	} else {
		var jeu Data
		jeu.Life = 10
		var filename string
		var filenameresult string
		fmt.Printf("Choose your difficulty : \n1 = Easy\n2 = Medium\n3 = Hard\n")
		fmt.Scanf("%s", &filename)
		switch filename {
		case "1":
			consoleReset()
			filenameresult = "words.txt"

		case "2":
			consoleReset()
			filenameresult = "words2.txt"

		case "3":
			consoleReset()
			filenameresult = "words3.txt"
		default:
			consoleReset()
			fmt.Printf("Error, default difficulty = 1\n\n")
			filenameresult = "Hangman-GO/words.txt"
		}
		fmt.Printf("Good luck, you have 10 attempts, you can stop the game with 'STOP'\n\n")
		hideWordStart(&jeu, randomWord(filenameresult))
	}
}
