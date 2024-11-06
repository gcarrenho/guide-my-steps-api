package translator

import "github.com/nicksnyder/go-i18n/v2/i18n"

type i18nRepo interface {
	LoadLanguage(languagePath string) (*i18n.MessageFile, error)
	GetLocalizer(language string) *i18n.Localizer
}
