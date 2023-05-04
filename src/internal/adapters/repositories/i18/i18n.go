package repositories

import (
	"encoding/json"
	"strings"

	models "project/guidemysteps/src/internal/adapters/repositories/open_street_map/models"

	usermodel "project/guidemysteps/src/internal/core/models"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type I18nRepo struct {
	bundle *i18n.Bundle
}

func NewI18nRepo(bundle *i18n.Bundle) I18nRepo {
	return I18nRepo{
		bundle: bundle,
	}
}

func (i I18nRepo) Translate(steps models.Steps, stepType string, nextInstruction string, user usermodel.User) string {
	languaje := strings.Split(user.Config.Language, "-")
	i.bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	i.bundle.MustLoadMessageFile("../internal/locales/" + languaje[0] + "/" + user.Config.Language + ".json")

	loc := i18n.NewLocalizer(i.bundle, "es")

	var i18nConfig i18n.LocalizeConfig

	translation := transalation(loc, steps, stepType, nextInstruction, i18nConfig, user)

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

func transalation(loc *i18n.Localizer, steps models.Steps, messageID string, nextInstruction string, i18nConfig i18n.LocalizeConfig, user usermodel.User) string {
	translation := ""
	streetName := steps.Name
	if streetName == "" {
		streetName = "sin nombre"
	}

	i18nConfig.TemplateData = map[string]interface{}{}

	if steps.Distance > 30 {
		messageID = strings.Split(messageID, "_pre")[0]
	}
	// Agregar el calculo de metros o pasos de acuerdo a la config del usuario.
	switch steps.Maneuver.Modifier {
	case "straight":
		translation = loc.MustLocalize(&i18n.LocalizeConfig{
			MessageID: messageID,
			TemplateData: map[string]interface{}{
				"RELATIVE_DIRECTION": loc.MustLocalize(&i18n.LocalizeConfig{
					MessageID:    steps.Maneuver.Modifier,
					TemplateData: map[string]interface{}{},
					PluralCount:  1,
				}),
				"STREET_NAMES": streetName,
				"UNITS":        user.Config.Unit,
				"COUNT_UNITS":  "x pasos",
				"NEXT_VERBAL":  nextInstruction,
				"LENGTH":       usermodel.GetLenghtStep(steps.Distance, user),
			},
		})
	default:
		if len(steps.Maneuver.Modifier) == 0 {
			steps.Maneuver.Modifier = pointCardinal(steps.Maneuver.BearingAfter)
		}

		translation = loc.MustLocalize(&i18n.LocalizeConfig{
			MessageID: messageID,
			TemplateData: map[string]interface{}{
				"RELATIVE_DIRECTION": loc.MustLocalize(&i18n.LocalizeConfig{
					MessageID:    steps.Maneuver.Modifier,
					TemplateData: map[string]interface{}{},
					PluralCount:  1,
				}),
				"STREET_NAMES": streetName,
				"UNITS":        user.Config.Unit,
				"COUNT_UNITS":  "x pasos", // se lo vamos a agregar desde flutter En x user.Config.units ....
				"NEXT_VERBAL":  nextInstruction,
				"LENGTH":       usermodel.GetLenghtStep(steps.Distance, user), //strconv.Itoa(int(steps.Distance / 0.762)),
			},
			PluralCount: 1,
		})
	}

	return translation
}

func PuntoCardinal(direction int) string {

	var cardinal string
	switch {
	case direction >= 338 || direction < 23:
		cardinal = "north"
	case direction >= 23 && direction < 68:
		cardinal = "north_east"
	case direction >= 68 && direction < 113:
		cardinal = "east"
	case direction >= 113 && direction < 158:
		cardinal = "south_east"
	case direction >= 158 && direction < 203:
		cardinal = "south"
	case direction >= 203 && direction < 248:
		cardinal = "south_west"
	case direction >= 248 && direction < 293:
		cardinal = "west"
	default:
		cardinal = "north_west"
	}

	return cardinal
}

var PathLenguaje = map[string]string{
	"es": "es/es-ES.json",
}
