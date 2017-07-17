package prompt

var InputHandler = defaultHandler

func defaultHandler(ac *ASCIICode, buffer *Buffer) {
	switch ac.Key {
	case Left:
		buffer.CursorLeft(1)
	case Right:
		buffer.CursorRight(1)
	case Backspace:
		buffer.DeleteBeforeCursor(1)
	case ControlI: // this is equivalent with TabKey.
		fallthrough
	case Tab:
		break
	case ControlT:
		break
		return
	case Up:
		break
	case Down:
		break
	default:
		break
	}
	return
}
