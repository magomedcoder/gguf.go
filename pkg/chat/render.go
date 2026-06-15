package chat

import (
	"fmt"
	"strings"

	"github.com/magomedcoder/gguf.go/pkg/format"
)

const (
	imStart = "<|im_start|>"
	imEnd   = "<|im_end|>"

	// Стандартные Qwen3 thinking-токены - при наличии metadata берутся из vocab
	defaultThinkingOpen  = " "
	defaultThinkingClose = " "
)

// Render форматирует диалог в ChatML/Qwen-стиле
// Использует messages и add_generation_prompt; наличие tokenizer.chat_template в метаданных проверяется через HasTemplateMeta, Jinja не парсится
func Render(meta format.Metadata, messages []Message, addGenerationPrompt bool) (string, error) {
	if !HasTemplateMeta(meta) {
		return "", fmt.Errorf("chat: tokenizer.chat_template не найден")
	}

	return renderChatML(messages, addGenerationPrompt), nil
}

// RenderFromReader рендерит prompt из GGUF reader
func RenderFromReader(r *format.Reader, messages []Message, addGenerationPrompt bool) (string, error) {
	return Render(r.Metadata, messages, addGenerationPrompt)
}

func renderChatML(messages []Message, addGenerationPrompt bool) string {
	var b strings.Builder

	for _, msg := range messages {
		switch msg.Role {
		case "system", "user", "assistant", "tool":
			writeBlock(&b, msg.Role, msg.Content)
		}
	}

	if addGenerationPrompt {
		writeAssistantPrompt(&b, true, nil)
	}

	return b.String()
}

func writeAssistantPrompt(b *strings.Builder, enableThinking bool, meta format.Metadata) {
	b.WriteString(imStart)
	b.WriteString("assistant\n")
	if !enableThinking {
		b.WriteString(EmptyThinkingBlock(meta))
	}
}

func writeBlock(b *strings.Builder, role, content string) {
	b.WriteString(imStart)
	b.WriteString(role)
	b.WriteByte('\n')
	b.WriteString(content)
	b.WriteString(imEnd)
	b.WriteByte('\n')
}
