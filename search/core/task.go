package core

type Task struct {
    Start int // The index of this chunk in the source text
    Pattern []byte // The search pattern
    Text []byte // The text of this chunk
}
