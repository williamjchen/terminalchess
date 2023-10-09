package server

type lobby struct {
	code string
	players []*player
	move bool
}