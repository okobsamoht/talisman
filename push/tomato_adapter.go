package push

import "github.com/okobsamoht/tomato/types"

type tomatoPushAdapter struct {
	validPushTypes []string
}

func newTomatoPush() *tomatoPushAdapter {
	t := &tomatoPushAdapter{
		validPushTypes: []string{"ios", "android"},
	}
	return t
}

func (t *tomatoPushAdapter) send(body types.M, installations types.S, pushStatus string) []types.M {
	return []types.M{}
}

func (t *tomatoPushAdapter) getValidPushTypes() []string {
	return t.validPushTypes
}
