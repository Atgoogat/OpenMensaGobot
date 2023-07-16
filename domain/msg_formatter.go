package domain

type TextFormatter interface {
	Format(text string) (string, error)
}
