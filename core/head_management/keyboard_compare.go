package head_management

import "github.com/televi-go/televi/models/render/results"

func KeyboardsAreSame(first, second results.ReplyKeyboardResult) bool {
	if len(first.Buttons) != len(second.Buttons) {
		return false
	}

	for i := 0; i < len(first.Buttons); i++ {
		if len(first.Buttons[i]) != len(second.Buttons[i]) {
			return false
		}

		for j := 0; j < len(first.Buttons[i]); j++ {
			if first.Buttons[i][j] != second.Buttons[i][j] {
				return false
			}
		}
	}

	return true
}
