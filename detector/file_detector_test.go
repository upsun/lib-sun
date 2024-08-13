package detector

import (
	"testing"

	entity "upsun.com/lib-upsun/entity"
)

const TEST_CONVERT_PATH = "../../tests/convert/"

func TestFind(t *testing.T) {
	want := make(map[string][]string)
	want[entity.PSH_ROUTE] = []string{TEST_CONVERT_PATH + ".platform/routes.yaml"}

	got, _ := FindConfig(TEST_CONVERT_PATH)

	if len(got[entity.PSH_ROUTE]) != 0 {
		if got[entity.PSH_ROUTE][0] != want[entity.PSH_ROUTE][0] {
			t.Errorf("Got %q, want %q", got[entity.PSH_ROUTE][0], want[entity.PSH_ROUTE][0])
		}
	} else {
		t.Errorf("No diretory found")
	}

}
