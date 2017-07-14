package prompt

var InputHandler = defaultHandler

func defaultHandler(ac *ASCIICode, buffer *Buffer, out *VT100Writer) {
	switch ac.Key {
	case ControlJ: // this is equivalent with Enter Key.
		fallthrough
	case Enter:
		out.EraseDown()
		out.WriteStr(buffer.Document().TextAfterCursor())

		out.WriteStr("\n>>> Your input: '")
		out.WriteStr(buffer.Text())
		out.WriteStr("' <<<\n")
		buffer = NewBuffer()
	case Left:
		l := buffer.CursorLeft(1)
		if l == 0 {
			return
		}
		out.EraseLine()
		out.EraseDown()
		after := buffer.Document().CurrentLine()
		out.WriteStr(after)
		out.CursorBackward(len(after) - buffer.CursorPosition)
	case Right:
		l := buffer.CursorRight(1)
		if l == 0 {
			return
		}

		out.CursorForward(l)
		out.WriteRaw(ac.ASCIICode)
		out.EraseDown()
		after := buffer.Document().TextAfterCursor()
		out.WriteStr(after)
	case Backspace:
		deleted := buffer.DeleteBeforeCursor(1)
		if deleted == "" {
			return
		}
		out.CursorBackward(1)
		out.EraseDown()

		after := buffer.Document().TextAfterCursor()
		out.WriteStr(after)
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
