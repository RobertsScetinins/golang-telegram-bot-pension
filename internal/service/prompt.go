package service

import (
	"fmt"
	"strings"
)

const inputPlaceholder = "{input}"

type PromptType string

const (
	PromptTypeFactCheck   PromptType = "fact_check"
	PromptTypeAnalyze     PromptType = "analyze_image"
	PromptTypeExtractText PromptType = "extract_text"
	PromptTypeCustom      PromptType = "custom"
	PromptTypeSummary     PromptType = "summary"
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
		UserPrompt: fmt.Sprintf("Утверждение: %s", inputPlaceholder),
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
		UserPrompt: inputPlaceholder,
	}

	pm.presets[PromptTypeExtractText] = &PromptPreset{
		Type: PromptTypeExtractText,
		SystemPrompt: `
			Сохраните исходное форматирование, разметку и язык.
			Если текст находится в таблице, сохраните табличную структуру.
			Выведите только извлеченный текст без дополнительных комментариев.
		`,
		UserPrompt: fmt.Sprintf("Извлеките текст из этого изображения: %s", inputPlaceholder),
	}

	pm.presets[PromptTypeSummary] = &PromptPreset{
		Type: PromptTypeSummary,
		SystemPrompt: `
			Ты — ассистент, создающий краткие и информативные саммари чат-переписок.

			Правила:
			- Отвечай только на русском языке
			- Выделяй основные темы обсуждения
			- Упоминай ключевые решения, вопросы или проблемы
			- Структурируй ответ для удобного чтения

			Формат вывода:
        	Краткое содержание переписки: 

		`,
		UserPrompt: fmt.Sprintf(`Вот переписка для анализа: %s Пожалуйста, создай саммари`, inputPlaceholder),
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
	if strings.TrimSpace(userInput) == "" {
		return strings.TrimSpace(p.SystemPrompt)
	}

	formattedUserPrompt := strings.ReplaceAll(p.UserPrompt, inputPlaceholder, userInput)

	return fmt.Sprintf("%s\n\n%s", p.SystemPrompt, formattedUserPrompt)
}
