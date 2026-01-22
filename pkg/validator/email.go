package validator

import (
	"regexp"
	"strings"
)

// Standard regex for email validation
var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

// EmailValidator defines the interface for email validation services.
type EmailValidator interface {
	ValidateFormat(email string) bool
	IsPersonalProvider(email string) bool
}

// emailValidator implements the EmailValidator interface.
type emailValidator struct{}

// NewEmailValidator creates a new EmailValidator.
func NewEmailValidator() EmailValidator {
	return &emailValidator{}
}

// ValidateFormat checks if the email has a valid format.
func (v *emailValidator) ValidateFormat(email string) bool {
	if email == "" || len(email) > 254 {
		return false
	}
	if strings.Count(email, "@") != 1 {
		return false
	}
	if strings.Contains(email, "..") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts[0]) > 64 {
		return false
	}
	return emailRegex.MatchString(email)
}

// IsPersonalProvider checks if the email domain belongs to a known personal provider.
func (v *emailValidator) IsPersonalProvider(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := strings.ToLower(parts[1])

	switch domain {
	// Major US providers
	case "gmail.com", "yahoo.com", "hotmail.com", "outlook.com", "aol.com", "icloud.com", "live.com", "msn.com":
		return true
	// Privacy-focused providers
	case "protonmail.com", "tutanota.com", "mailbox.org", "startmail.com", "runbox.com", "posteo.de", "kolabnow.com", "disroot.org", "riseup.net":
		return true
	// Other popular providers
	case "mail.com", "zoho.com", "fastmail.com", "yandex.com":
		return true
	// Temporary/disposable email providers
	case "guerrillamail.com", "10minutemail.com", "tempmail.org", "mailinator.com":
		return true
	// International & country-specific variants
	case "yahoo.co.uk", "yahoo.ca", "yahoo.in", "yahoo.com.br", "yahoo.co.jp", "yahoo.de", "yahoo.fr", "yahoo.it", "yahoo.es",
		"hotmail.co.uk", "hotmail.fr", "hotmail.de", "hotmail.it", "hotmail.es", "hotmail.com.br", "outlook.com.br", "outlook.co.uk",
		"outlook.fr", "outlook.de", "outlook.it", "outlook.es", "rediffmail.com", "sify.com", "in.com", "163.com", "126.com", "qq.com",
		"sina.com", "sohu.com", "yeah.net", "tom.com", "naver.com", "daum.net", "hanmail.net", "nate.com", "mail.ru", "rambler.ru",
		"bk.ru", "list.ru", "inbox.ru", "t-online.de", "gmx.de", "gmx.com", "web.de", "freenet.de", "arcor.de", "orange.fr",
		"laposte.net", "free.fr", "wanadoo.fr", "sfr.fr", "alice.fr", "libero.it", "virgilio.it", "alice.it", "tin.it", "tiscali.it",
		"fastwebnet.it", "terra.com.br", "uol.com.br", "bol.com.br", "globo.com", "ig.com.br", "r7.com", "oi.com.br", "terra.es",
		"ya.com", "telefonica.net", "btinternet.com", "sky.com", "virgin.net", "tiscali.co.uk", "sympatico.ca", "rogers.com",
		"bell.net", "shaw.ca", "bigpond.com", "optusnet.com.au", "iinet.net.au", "westnet.com.au", "so-net.ne.jp", "nifty.com",
		"biglobe.ne.jp", "ocn.ne.jp":
		return true
	default:
		return false
	}
}
