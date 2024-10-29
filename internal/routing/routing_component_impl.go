package routing

import "context"

var _ RoutingComponent = (*routingComponentImpl)(nil)

type routingComponentImpl struct {
	routingProvider routingProvider
}

func NewRoutingComponentImpl(routingProvider routingProvider) *routingComponentImpl {
	return &routingComponentImpl{
		routingProvider: routingProvider,
	}
}

func (o routingComponentImpl) GetRouting(ctx context.Context, myStepResquest RoutesRequest) (MySteps, error) {
	return o.routingProvider.GetRouting(ctx, myStepResquest)
}

/*
func (s *RoutingComponentImpl) TranslateSteps(steps Steps, stepType, nextInstruction string, user appuser.User) string {
	languageCode := strings.Split(user.Config.Language, "-")[0]
	localizer := s.localizationRepo.GetLocalizer(languageCode)

	return buildTranslation(localizer, steps, stepType, nextInstruction, user)
}

func buildTranslation(localizer *i18n.Localizer, steps routing.Steps, messageID, nextInstruction string, user appuser.User) string {
	// Lógica para construir la instrucción de traducción
	// ...
	translation := ""
	streetName := steps.Name
	if streetName == "" {
		streetName = "sin nombre"
	}

	//i18nConfig.TemplateData = map[string]interface{}{}

	if steps.Distance > 30 {
		messageID = strings.Split(messageID, "_pre")[0]
	}
	// Agregar el calculo de metros o pasos de acuerdo a la config del usuario.
	switch steps.Maneuver.Modifier {
	case "straight":
		translation = localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: messageID,
			TemplateData: map[string]interface{}{
				"RELATIVE_DIRECTION": localizer.MustLocalize(&i18n.LocalizeConfig{
					MessageID:    steps.Maneuver.Modifier,
					TemplateData: map[string]interface{}{},
					PluralCount:  1,
				}),
				"STREET_NAMES": streetName,
				"UNITS":        user.Config.Unit,
				"COUNT_UNITS":  "x pasos",
				"NEXT_VERBAL":  nextInstruction,
				"LENGTH":       user.GetLengthStep(steps.Distance, user),
			},
		})
	default:
		//if len(steps.Maneuver.Modifier) == 0 {
		steps.Maneuver.Modifier = pointCardinal(steps.Maneuver.BearingAfter)
		//}

		translation = localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: messageID,
			TemplateData: map[string]interface{}{
				"RELATIVE_DIRECTION": localizer.MustLocalize(&i18n.LocalizeConfig{
					MessageID:    steps.Maneuver.Modifier,
					TemplateData: map[string]interface{}{},
					PluralCount:  1,
				}),
				"STREET_NAMES": streetName,
				"UNITS":        user.Config.Unit,
				"COUNT_UNITS":  "x pasos", // se lo vamos a agregar desde flutter En x user.Config.units ....
				"NEXT_VERBAL":  nextInstruction,
				"LENGTH":       user.GetLengthStep(steps.Distance, user), //strconv.Itoa(int(steps.Distance / 0.762)),
			},
			PluralCount: 1,
		})
	}

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

}*/

/*

func buildTranslation(localizer *i18n.Localizer, step TranslatableStep, messageID, nextInstruction string, user usermodel.User) string {
	streetName := step.GetStreetName()
	if streetName == "" {
		streetName = "sin nombre"
	}

	// Ejemplo simplificado de la construcción de una instrucción traducida.
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: messageID,
		TemplateData: map[string]interface{}{
			"STREET_NAMES": streetName,
			"UNITS":        user.Config.Unit,
			"NEXT_VERBAL":  nextInstruction,
			"LENGTH":       usermodel.GetLengthStep(step.GetDistance(), user),
		},
	})
}
*/
