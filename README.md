# Hangman - GO

Ce programme permet de faire le jeu du pendu en GO utilisant la version go 1.22.2

## Fonctionnalités

- Choix de la difficulté : Trois niveaux disponibles avec words.txt, words2.txt et words3.txt avec des fichiers de mots différents.
- Sauvegarde et chargement : La partie peut être sauvegardée et chargée à tout moment.
- ASCII art : Visuel propre avec l'affichage du pendu
- UI : Simple et intuitif 

## Utilisation

### Lancer une partie

- Pour lancer une partie, appeler la fonction GameStart du module hangman via un fichier main.go (Déja present dans le projet)
- Donc pour le lancer il vous suffit d'appeler ```bash go run main/main.go``` du projet

### Sauvegarder et charger une partie

- Sauvegarder en entrant `STOP`
- Charger une partie avec :
```bash
go run main/main.go --startWith game.txt
```
