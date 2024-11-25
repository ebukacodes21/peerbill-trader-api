package mail

import (
	"peerbill-trader-api/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	config, err := utils.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSender, config.EmailAddress, config.EmailPassword)
	subject := "a test mail"
	content := `
	<h1>hello there</h1>
	<p>mail from peerbill</p>
	`
	to := []string{"georgeokafo1@gmail.com"}
	files := []string{"../README.md"}

	err = sender.SendMail(subject, content, to, nil, nil, files)
	require.NoError(t, err)
}
