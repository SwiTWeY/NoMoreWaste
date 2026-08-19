package mailer

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/SwiTWeY/NoMoreWaste/api/config"
	"github.com/SwiTWeY/NoMoreWaste/api/models"
)

//go:embed rappel_renouvellement.html
var rappelTemplate string

func EnvoyerRappel(cfg config.Config, r models.RappelAdhesion) error {
	tmpl, err := template.New("rappel").Parse(rappelTemplate)
	if err != nil {
		return err
	}

	data := struct {
		Prenom  string
		Nom     string
		DateFin string
	}{
		Prenom:  r.Prenom,
		Nom:     r.Nom,
		DateFin: r.DateFin.Format("02/01/2006"),
	}

	var corps bytes.Buffer
	if err := tmpl.Execute(&corps, data); err != nil {
		return err
	}

	message := fmt.Sprintf("From: NO MORE WASTE <%s>\r\n"+
		"To: %s\r\n"+
		"Subject: Rappel : renouvellement de votre adhesion\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s", cfg.SMTPUser, r.Email, corps.String())

	addr := cfg.SMTPHost + ":" + cfg.SMTPPort
	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPHost)

	return smtp.SendMail(addr, auth, cfg.SMTPUser, []string{r.Email}, []byte(message))
}
