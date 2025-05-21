package i18n

import (
	"errors"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/obumax/pet-password-generator/internal/generator"
)

// LocalizeError преобразует ошибку из generator в читабельное сообщение
// LocalizeError converts the error from generator into a human-readable message
func LocalizeError(localizer *i18n.Localizer, err error, data map[string]interface{}) string {
	var msgID string
	switch {
	case errors.Is(err, generator.ErrLengthOutOfRange):
		msgID = "length_out_of_range"
	case errors.Is(err, generator.ErrNoCategorySelected):
		msgID = "no_category_selected"
	default:
		// немашинная ошибка
		// non-machine error
		return err.Error()
	}
	str, _ := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    msgID,
		TemplateData: data,
	})
	return str
}
