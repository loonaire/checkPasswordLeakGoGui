package main

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	app := app.New()
	window := app.NewWindow("GoCheckPasswordLeak")

	// chargement des widget
	entryLabel := widget.NewLabel("Mot de passe à tester:")
	entryPassword := widget.NewPasswordEntry()

	resultLabel := widget.NewLabel("")

	buttonValidation := widget.NewButton("Tester le mot de passe", func() {
		if entryPassword.Text != "" {
			// hash le mot de passe
			passwordHash := HashString(entryPassword.Text)

			// on récupère les hashs qui possèdent les mêmes 5 premiers caractère que notre hash
			strHashPossible, err := GetHtmlFromUrl("https://api.pwnedpasswords.com/range/" + passwordHash[:5])
			if err != nil {
				// si erreur lors de la récupération des infos sur le site haveibeenpwned
				resultLabel.SetText("Erreur lors de la récupération des informations:" + err.Error())
			} else {
				var passwordLeaked bool = false
				// on split les infos une première fois pour charger chaque possibilité
				tabPossibilities := strings.Split(strHashPossible, "\n")

				for _, elt := range tabPossibilities {
					// on resplite pour avoir le hash à tester ainsi que le nombre de fois ou il a leak si il a leak
					tabElt := strings.Split(elt, ":")
					if tabElt[0] == passwordHash[5:] {
						resultLabel.SetText("Le mot de passe à fuité dans " + tabElt[1] + " leak")
						passwordLeaked = true
						break
					}
				}
				if !passwordLeaked {
					resultLabel.SetText("Le mot de passe n'a pas fuité")
				}

			}

		}
	})

	// mise en forme
	entryContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(2), entryLabel, entryPassword)

	globalContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(1), entryContainer, buttonValidation, resultLabel)

	window.SetContent(globalContainer)

	window.Show()

	app.Run()

}

func HashString(str string) string {
	h := sha1.New()
	h.Write([]byte(str))

	hash := fmt.Sprintf("%X", h.Sum(nil))

	return hash

}

func GetHtmlFromUrl(url string) (string, error) {

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Erreur ", err)
		return "", errors.New("Problème de connexion internet ou erreur dans l'url")
	}

	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Erreur", err)
		return "", errors.New("Erreur lors de la lecture ")
	}

	return string(html), nil
}
