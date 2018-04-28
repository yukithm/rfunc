package rfuncs

func MakeTextClipboardContent(text string) *ClipboardContent {
	return &ClipboardContent{
		Type: ClipboardType_TEXT,
		Content: &ClipboardContent_Text{
			Text: text,
		},
	}
}

func (m *PasteRequest) Acceptable(t ClipboardType) bool {
	accepts := m.GetAccepts()
	if accepts == nil || len(accepts) == 0 {
		return false
	}

	for _, accept := range accepts {
		if accept == t {
			return true
		}
	}

	return false
}
