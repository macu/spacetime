package treetime

type NodeHeader struct {
	ID    uint      `json:"id"`
	Class string    `json:"class"`
	Key   *string   `json:"key"`
	Title string    `json:"title"`
	Body  string    `json:"body"`
	Tags  []NodeTag `json:"tags"`
}
