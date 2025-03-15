package model

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func NewTask(title string) *Task {
	return &Task{Title: title, Done: false}
}

func (t *Task) MarkDone() {
	t.Done = true
}

func (t *Task) MarkUndone() {
	t.Done = false
}

func (t *Task) Edit(title string) {
	t.Title = title
}
