package service

import (
	"fmt"
	"strings"
)

type PromptType string

const (
	PromptTypeFactCheck   PromptType = "fact_check"
	PromptTypeAnalyze     PromptType = "analyze_image"
	PromptTypeExtractText PromptType = "extract_text"
	PromptTypeCustom      PromptType = "custom"
)

type PromptPreset struct {
	Type         PromptType
	SystemPrompt string
	UserPrompt   string
	Options      PromptOptions
}

type PromptOptions struct {
	UseGoogleSearch bool
	Temperature     float64
	MaxTokens       int
	ResponseFormat  string
}

type PromptManager struct {
	presets map[PromptType]*PromptPreset
}

func NewPromptManager() *PromptManager {
	pm := &PromptManager{
		presets: make(map[PromptType]*PromptPreset),
	}

	pm.registerDefaultPresets()

	return pm
}

func (pm *PromptManager) registerDefaultPresets() {
	pm.presets[PromptTypeFactCheck] = &PromptPreset{
		Type: PromptTypeFactCheck,
		SystemPrompt: `
			Вы – помощник по проверке фактов.

			Правила:
			- Сохраняйте язык утверждения.
			- Отвечайте кратко.
			- Всегда отвечайте на русском языке.
			- Если утверждение касается событий, дат, статистики или проверяемых фактов, используйте интернет для поиска свежей информации.
			- Если утверждение субъективное или личное мнение – не гуглите.

			Формат вывода:
			Утверждение: <claim>
			Оценка: Факт | Ложь | Вводящее в заблуждение | Мнение | Неверифицируемо
			Объяснение: <максимум 40 слов>
			Источники: <список ссылок, если использовалась проверка в интернете>
		`,
		UserPrompt: "Утверждение: {input}",
	}

	pm.presets[PromptTypeAnalyze] = &PromptPreset{
		Type: PromptTypeAnalyze,
		SystemPrompt: `
		Вы — эксперт в области анализа.

		Подробно опишите изображение, включая:
		- Основные объекты/предметы
		- Цвета и композиция
		- Текст, видимый на изображении
		- Любые действия или взаимодействия
		- Общий контекст и настроение
		`,
		UserPrompt: "{input}",
	}

	pm.presets[PromptTypeExtractText] = &PromptPreset{
		Type: PromptTypeExtractText,
		SystemPrompt: `
			Сохраните исходное форматирование, разметку и язык.
			Если текст находится в таблице, сохраните табличную структуру.
			Выведите только извлеченный текст без дополнительных комментариев.
		`,
		UserPrompt: "Извлеките текст из этого изображения: {input}",
	}
}

func (pm *PromptManager) GetPromptPreset(promptType PromptType) (*PromptPreset, error) {
	preset, exists := pm.presets[promptType]
	if !exists {
		return nil, fmt.Errorf("prompt preset not found: %s", promptType)
	}

	return preset, nil
}

// FormatPrompt combines system prompt and user input
func (p *PromptPreset) FormatPrompt(userInput string) string {
	formattedUserPrompt := strings.ReplaceAll(p.UserPrompt, "{input}", userInput)

	return fmt.Sprintf("%s\n\n%s", p.SystemPrompt, formattedUserPrompt)
}
