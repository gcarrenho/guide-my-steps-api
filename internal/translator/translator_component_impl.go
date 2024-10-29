package translator

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var _ Translator = (*TranslationService)(nil)

type TranslationService struct {
	localizationRepo *i18nRepoImpl
}

func NewTranslationService(repo *i18nRepoImpl) *TranslationService {
	return &TranslationService{localizationRepo: repo}
}

func (s *TranslationService) TranslateStep(translatorParam TranslatorStep) string {
	languageCode := strings.Split(translatorParam.Config.Language, "-")[0]

	langFilePath := filepath.Join("../../../internal/translator/locales", languageCode, translatorParam.Config.Language+".json")
	absolutePath, err := filepath.Abs(langFilePath)
	if err != nil {
		fmt.Println("Errooooor ", err)
		//return fmt.Println("Errooooor ", err)
	}

	_, err = s.localizationRepo.LoadLanguage(absolutePath)
	if err != nil {
		fmt.Println("Errooooor ", err)
	}

	localizer := s.localizationRepo.GetLocalizer(languageCode)

	//i.bundle.MustLoadMessageFile("../internal/locales/" + languageCode + "/" + translatorParam.Config.Language + ".json")

	return buildTranslation(localizer, translatorParam)
}

func buildTranslation(localizer *i18n.Localizer, translatorParam TranslatorStep) string {
	// Lógica para construir la instrucción de traducción
	// ...
	translation := ""
	streetName := translatorParam.Name
	if streetName == "" {
		streetName = "sin nombre"
	}

	//i18nConfig.TemplateData = map[string]interface{}{}

	if translatorParam.Distance > 30 {
		translatorParam.MessageID = strings.Split(translatorParam.MessageID, "_pre")[0]
	}

	relativeDirection := func(modifier string, bearingAfter int) string {
		if modifier == "straight" {
			return modifier
		}
		// Si no hay modifier, calculamos el punto cardinal en función del bearing
		return pointCardinal(bearingAfter)
	}

	templateData := map[string]interface{}{
		"RELATIVE_DIRECTION": localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID:    relativeDirection(translatorParam.Modifier, translatorParam.BearingAfter),
			TemplateData: map[string]interface{}{},
			PluralCount:  1,
		}),
		"STREET_NAMES": streetName,
		"UNITS":        translatorParam.Config.Unit,
		"COUNT_UNITS":  "x pasos", // Esto se puede ajustar según la configuración del usuario en Flutter
		"NEXT_VERBAL":  translatorParam.NextInstruction,
		"LENGTH":       translatorParam.LengthStep, //user.GetLengthStep(translatorParam.Distance, user),
	}

	translation = localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    translatorParam.StepType,
		TemplateData: templateData,
		PluralCount:  1,
	})

	return translation
}

func pointCardinal(bearing int) string {
	var cardinal string
	switch {
	case bearing >= 338 || bearing < 23:
		cardinal = "north"
	case bearing >= 23 && bearing < 68:
		cardinal = "north_east"
	case bearing >= 68 && bearing < 113:
		cardinal = "east"
	case bearing >= 113 && bearing < 158:
		cardinal = "south_east"
	case bearing >= 158 && bearing < 203:
		cardinal = "south"
	case bearing >= 203 && bearing < 248:
		cardinal = "south_west"
	case bearing >= 248 && bearing < 293:
		cardinal = "west"
	default:
		cardinal = "north_west"
	}

	return cardinal

}
