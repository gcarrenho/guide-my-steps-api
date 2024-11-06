package translator

import (
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var _ i18nRepo = (*i18nRepoImpl)(nil)

type i18nRepoImpl struct {
	bundle *i18n.Bundle
}

func NewI18nRepo(bundle *i18n.Bundle) *i18nRepoImpl {
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	return &i18nRepoImpl{bundle: bundle}
}

func (r *i18nRepoImpl) LoadLanguage(languagePath string) (*i18n.MessageFile, error) {
	return r.bundle.LoadMessageFile(languagePath)
}

func (r *i18nRepoImpl) GetLocalizer(language string) *i18n.Localizer {
	return i18n.NewLocalizer(r.bundle, language)
}
