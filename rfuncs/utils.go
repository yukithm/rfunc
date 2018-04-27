package rfuncs

func MakeTextClipboardContent(text string) *ClipboardContent {
	return &ClipboardContent{
		Type: ClipboardType_TEXT,
		Content: &ClipboardContent_Text{
			Text: text,
		},
	}
}
