package valueobject

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/mjiee/gokit"

	"github.com/mjiee/world-news/backend/pkg/locale"
	"github.com/mjiee/world-news/backend/pkg/ttsai"
)

var (
	systemKey     = "system"
	classifyKey   = "classify"
	scriptKey     = "script"
	announcerKey  = "announcer"
	emotionKey    = "emotion"
	speechRateKey = "speechRate"
	volumeKey     = "volume"
	scriptJsonKey = "scriptJson"
	mergeKey      = "merge"
	approvalKey   = "approval"
)

// buildKey builds the key
func buildKey(key, language string) string {
	return fmt.Sprintf("%s_%s", key, language)
}

// Prompt locale
var promptLocale = map[string]string{
	buildKey(systemKey, locale.En):     "You are a podcast script generation assistant.",
	buildKey(systemKey, locale.Zh):     "你是一位中文播客脚本生成助手。",
	buildKey(classifyKey, locale.En):   "Please analyze the news content and determine which podcast category it is most suitable for, returning the index of the category (single choice):",
	buildKey(classifyKey, locale.Zh):   "请分析新闻内容，判断它最适合哪种播客分类，只返回分类索引（单选）：",
	buildKey(scriptKey, locale.En):     "Convert podcast content into tts scripts.",
	buildKey(scriptKey, locale.Zh):     "将播客内容转换为tts脚本。",
	buildKey(announcerKey, locale.En):  "Radio announcer: ",
	buildKey(announcerKey, locale.Zh):  "播音员: ",
	buildKey(emotionKey, locale.En):    `Emotion: ["happy", "sad", "angry", "surprised", "fear", "hate", "excited", "coldness", "neutral", "depressed", "lovey-dovey", "shy", "comfort", "tension", "tender", "storytelling", "radio", "magnetic", "advertising", "vocal - fry", "ASMR", "news", "entertainment"]`,
	buildKey(emotionKey, locale.Zh):    `情绪: ["happy", "sad", "angry", "surprised", "fear", "hate", "excited", "coldness", "neutral", "depressed", "lovey-dovey", "shy", "comfort", "tension", "tender", "storytelling", "radio", "magnetic", "advertising", "vocal - fry", "ASMR", "news", "entertainment"]`,
	buildKey(speechRateKey, locale.En): "Speech rate: [0,2]",
	buildKey(speechRateKey, locale.Zh): "语速: [0,2]",
	buildKey(volumeKey, locale.En):     "Volume: [0,100]",
	buildKey(volumeKey, locale.Zh):     "音量: [0,100]",
	buildKey(scriptJsonKey, locale.En): "Only the standard json list is required to be output. Example: ",
	buildKey(scriptJsonKey, locale.Zh): "要求只输出标准json列表，示例：",
	buildKey(mergeKey, locale.En):      "Please merge the following podcast content into a single podcast script.",
	buildKey(mergeKey, locale.Zh):      "请合并多篇播客内容，保留所有文本的核心信息，使其成为一篇完整的播客文案。",
	buildKey(approvalKey, locale.En):   `Please reply with "yes" or "no" and provide the reason.`,
	buildKey(approvalKey, locale.Zh):   `请回复“yes”或“no”，并给出理由。`,
}

// getDefaultPrompt returns the default prompt
func getDefaultPrompt(key, language string) string {
	return promptLocale[buildKey(key, language)]
}

// PodcastScriptPrompt represents the prompt for podcast script generation
type PodcastScriptPrompt struct {
	SystemPrompt   string         `json:"systemPrompt"`
	ApprovalPrompt string         `json:"approvalPrompt"`
	RewritePrompt  string         `json:"rewritePrompt"`
	MergePrompt    string         `json:"mergePrompt"`
	ClassifyPrompt string         `json:"classifyPrompt"`
	StylizePrompts []*StylePrompt `json:"stylizePrompts"`
	ScriptPrompt   string         `json:"scriptPrompt"`
}

// StylePrompt represents the style prompt for podcast script generation
type StylePrompt struct {
	Style  string `json:"style"`
	Prompt string `json:"prompt"`
}

// BuildSystemPrompt returns the system prompt
func (p *PodcastScriptPrompt) BuildSystemPrompt(language string) string {
	if p.SystemPrompt != "" {
		return p.SystemPrompt
	}

	return getDefaultPrompt(systemKey, language)
}

// BuildApprovalPrompt returns the audit prompt
func (p *PodcastScriptPrompt) BuildApprovalPrompt(language string) string {
	return fmt.Sprintf("%s\n%s", p.ApprovalPrompt, getDefaultPrompt(approvalKey, language))
}

// BuildMergePrompt returns the merge prompt
func (p *PodcastScriptPrompt) BuildMergePrompt(language string) string {
	prompt := p.MergePrompt

	if prompt == "" {
		prompt = getDefaultPrompt(mergeKey, language)
	}

	return prompt
}

// BuildClassifyPrompt returns the classify prompt
func (p *PodcastScriptPrompt) BuildClassifyPrompt(language string) string {
	prompt := p.ClassifyPrompt

	if prompt == "" {
		prompt = getDefaultPrompt(classifyKey, language)
	}

	for idx, style := range p.StylizePrompts {
		prompt = fmt.Sprintf("%s\n%d. %s;\n", prompt, idx+1, style.Style)
	}

	return prompt
}

// GetStylePrompt returns the style
func (p *PodcastScriptPrompt) GetStylePrompt(output string) *StylePrompt {
	for _, item := range strings.Split(output, "\n") {
		if item != "" {
			output = item
			break
		}
	}

	for idx, style := range p.StylizePrompts {
		if strings.Contains(output, strconv.Itoa(idx+1)) {
			return style
		}
	}

	return nil
}

// VerifyApprovalResult verifies the approval result
func (p *PodcastScriptPrompt) VerifyApprovalResult(result string) StageStatus {
	for _, item := range strings.Split(result, "\n") {
		if item != "" {
			result = item
			break
		}
	}

	if strings.Contains(strings.ToLower(result), "yes") {
		return StageStatusCompleted
	}

	return StageStatusFailed
}

// ExtractScripts extracts the scripts
func (p *PodcastScriptPrompt) ExtractScripts(result string) []*ttsai.TtsScript {
	var (
		scripts []*ttsai.TtsScript
		start   = strings.Index(result, "[")
	)

	if start == -1 {
		return scripts
	}

	jsonStr := extractFromPosition(result, start, '[', ']')
	if jsonStr == "" {
		return scripts
	}

	if err := json.Unmarshal([]byte(jsonStr), &scripts); err != nil {
		return nil
	}

	return scripts
}

// extractFromPosition extracts the text from the given position
func extractFromPosition(text string, start int, open, close byte) string {
	depth := 0
	inString := false
	escape := false

	for i := start; i < len(text); i++ {
		c := text[i]

		if escape {
			escape = false
			continue
		}
		if c == '\\' {
			escape = true
			continue
		}

		if c == '"' {
			inString = !inString
			continue
		}
		if inString {
			continue
		}

		if c == open {
			depth++
		} else if c == close {
			depth--
			if depth == 0 {
				jsonStr := text[start : i+1]
				if json.Valid([]byte(jsonStr)) {
					return jsonStr
				}
				return ""
			}
		}
	}

	return ""
}

// BuildScriptPrompt builds the script prompt
func BuildScriptPrompt(language string, voices []*ttsai.Voice) string {
	var (
		prompt  = getDefaultPrompt(scriptKey, language)
		scripts = []*ttsai.TtsScript{}
	)

	for idx, voice := range voices {
		if idx == 0 {
			voice.Role = "MC"
		} else {
			voice.Role = "Co-Host"
		}
		scripts = append(scripts, &ttsai.TtsScript{
			Speaker:    voice.Id,
			Content:    fmt.Sprintf("This is the %d paragraph", idx),
			Emotion:    "news",
			SpeechRate: 1.0,
			Volume:     50,
		})
	}

	prompt = fmt.Sprintf("%s\n%s%s", prompt, getDefaultPrompt(announcerKey, language), gokit.MarshalSafe(voices))
	prompt = fmt.Sprintf("%s\n%s", prompt, getDefaultPrompt(emotionKey, language))
	prompt = fmt.Sprintf("%s\n%s", prompt, getDefaultPrompt(speechRateKey, language))
	prompt = fmt.Sprintf("%s\n%s", prompt, getDefaultPrompt(volumeKey, language))
	prompt = fmt.Sprintf("%s\n%s%s", prompt, getDefaultPrompt(scriptJsonKey, language), gokit.MarshalSafe(scripts))

	return prompt
}
