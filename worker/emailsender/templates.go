package emailsender

import "io/ioutil"

const userInviteTextTmpl = "worker/emailsender/templates/invite.txt"
const userInviteHTMLTmpl = "worker/emailsender/templates/invite.html"
const passwordResetTextTmpl = "worker/emailsender/templates/passwordreset.txt"
const passwordResetHTMLTmpl = "worker/emailsender/templates/passwordreset.html"

// MailTemplates contain all the templates in memory required for enqueuing
// mail jobs.
type MailTemplates struct {
	PasswordReset *MailTemplate
	UserInvite    *MailTemplate
}

// MailTemplate is a set of templates for emails, specifically for the MIME
// type text/plain and text/html.
type MailTemplate struct {
	TextTemplate string
	HTMLTemplate string
}

// NewMailTemplates returns a new instance of MailTemplates.
func NewMailTemplates() (*MailTemplates, error) {
	uit, err := NewMailTemplate(userInviteTextTmpl, userInviteHTMLTmpl)
	if err != nil {
		return nil, err
	}
	prt, err := NewMailTemplate(passwordResetTextTmpl, passwordResetHTMLTmpl)
	if err != nil {
		return nil, err
	}
	return &MailTemplates{UserInvite: uit, PasswordReset: prt}, nil
}

// NewEmailTemplate returns a new instance of MailTemplate.
func NewMailTemplate(textPath, htmlPath string) (*MailTemplate, error) {
	textTemplate, err := ioutil.ReadFile(textPath)
	if err != nil {
		return nil, err
	}
	htmlTemplate, err := ioutil.ReadFile(htmlPath)
	if err != nil {
		return nil, err
	}
	return &MailTemplate{
		TextTemplate: string(textTemplate),
		HTMLTemplate: string(htmlTemplate),
	}, nil
}
