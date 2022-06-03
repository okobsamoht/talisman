package push

import "github.com/okobsamoht/talisman/types"

type talismanPushAdapter struct {
	validPushTypes []string
}

func newTalismanPush() *talismanPushAdapter {
	t := &talismanPushAdapter{
		validPushTypes: []string{"ios", "android"},
	}
	return t
}

func (t *talismanPushAdapter) send(body types.M, installations types.S, pushStatus string) []types.M {
	return []types.M{}
}

func (t *talismanPushAdapter) getValidPushTypes() []string {
	return t.validPushTypes
}
