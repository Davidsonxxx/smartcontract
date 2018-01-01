package dialogManager

import (
	"gitlab.com/gameraccoon/telegram-accountant-bot/processing"
)

type TextProcessorFunc func(string, *processing.ProcessData) bool

type TextProcessorsMap map[string]TextProcessorFunc

type TextInputProcessorManager struct {
	Processors TextProcessorsMap
}

func (textProcessorManager *TextInputProcessorManager) processText(data *processing.ProcessData) bool {
	textProcessor := data.Static.GetUserStateTextProcessor(data.UserId)
	if textProcessor != nil {
		processor, ok := textProcessorManager.Processors[textProcessor.ProcessorId]
		if ok {
			return processor(textProcessor.AdditionalId, data)
		}
	}
	return false
}