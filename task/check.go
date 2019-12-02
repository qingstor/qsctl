package task

func (t *IsDestinationObjectExistTask) new() {}
func (t *IsDestinationObjectExistTask) run() {
	t.SetResult(t.GetDestinationObject() != nil)
}

func (t *IsSizeEqualTask) new() {}
func (t *IsSizeEqualTask) run() {
	t.SetResult(t.GetSourceObject().Size == t.GetDestinationObject().Size)
}

func (t *IsUpdateAtGreaterTask) new() {}
func (t *IsUpdateAtGreaterTask) run() {
	t.SetResult(
		t.GetSourceObject().UpdatedAt.After(t.GetDestinationObject().UpdatedAt),
	)
}
