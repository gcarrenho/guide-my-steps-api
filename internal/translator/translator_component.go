package translator

type Translator interface {
	TranslateStep(translatorParam TranslatorStep) string
}
