package interfaces

// Chat interface represents the unified interface for different chat service providers
type Chat interface {
	SendMessage(sender string, recipient string, message string) error
}
