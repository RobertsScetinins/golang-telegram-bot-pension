package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/helpers"
	"github.com/go-telegram/bot/models"
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

	pm.presets[PromptTypeCustom] = &PromptPreset{
		Type: PromptTypeCustom,
		SystemPrompt: `
			Ты — универсальный ИИ-помощник, который адаптируется под любые задачи пользователя.

			ТВОЙ ХАРАКТЕР:
			- Дружелюбный, но не приторный
			- Помогаешь по делу, без воды
			- Если не знаешь — честно говоришь, что не знаешь
			- Можешь пошутить, но только если это уместно
			- Адаптируешься под стиль общения пользователя

			ТВОИ ВОЗМОЖНОСТИ:
			- Отвечаешь на вопросы любой тематики
			- Помогаешь с написанием текстов, кода, писем
			- Объясняешь сложные вещи простым языком
			- Даёшь советы и рекомендации
			- Можешь выступать в роли эксперта в разных областях

			ПРАВИЛА ОТВЕТОВ:
			- Отвечай на том же языке, на котором написан запрос
			- Если просят что-то конкретное — делай именно это
			- Если просят совета — объясни плюсы и минусы
			- Если просят объяснить — начни с простого, потом углубляй
			- Для кода: давай готовые примеры с комментариями
			- Для текстов: предлагай варианты и улучшения

			ЧЕГО НЕ ДЕЛАТЬ:
			- Не выдумывай факты
			- Не давай медицинских/юридических консультаций (предупреди, что это только для информации)
			- Не поддерживай вредные или опасные действия
			- Не будь слишком формальным, если не просят
		`,
		UserPrompt: fmt.Sprintf(`Запрос пользователя: %s`, inputPlaceholder),
	}

	pm.presets[PromptTypeSummary] = &PromptPreset{
		Type: PromptTypeSummary,
		SystemPrompt: `
			Ты — дерзкий и остроумный рассказчик, который делает краткие, но максимально сочные саммари чат-переписок.
	
			Твоя задача — не просто пересказать, а превратить переписку в живой, местами саркастичный и забавный рассказ без цензуры (разрешён мат, если он уместен по контексту).
	
			Правила:
			- Пиши в неформальном, разговорном стиле, как будто пересказываешь другу
			- Сохраняй суть диалога, но добавляй юмор, иронию и немного драматизации
			- Чётко отмечай: кто что сказал, кто с кем спорил или соглашался
			- Подсвечивай конфликты, странные моменты и тупики в разговоре
			- Можешь слегка утрировать эмоции, если это делает пересказ смешнее
			- Не превращай текст в список — это должен быть цельный рассказ
			- Не будь слишком длинным, но и не теряй важные детали
	
			Формат:
			Один связный абзац или несколько коротких абзацев, читается легко и с юмором.
		`,
		UserPrompt: fmt.Sprintf(`Вот переписка для анализа: %s Сделай максимально живой и забавный пересказ`, inputPlaceholder),
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

func BuildPromptInput(message *models.Message, hasMedia bool) (string, error) {
	var userInput string

	if hasMedia && message.Caption != "" {
		userInput = message.Caption
	} else {
		userInput = message.Text
	}

	userComment, hasArgs := helpers.GetCommandArgs(userInput)
	var replyText string
	if message.ReplyToMessage != nil {
		reply := message.ReplyToMessage

		if reply.Text != "" {
			replyText = reply.Text
		} else if reply.Caption != "" {
			replyText = reply.Caption
		}
	}

	if !hasArgs && replyText == "" && !hasMedia {
		return "", errors.New("empty prompt")
	}

	var builder strings.Builder

	if hasArgs {
		builder.WriteString(userComment)
	}

	if replyText != "" {
		if builder.Len() > 0 {
			builder.WriteString("\n\n")
		}
		builder.WriteString("Контекст сообщения:\n")
		builder.WriteString(replyText)
	}

	return builder.String(), nil

}
