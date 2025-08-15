package carta

import (
	"reflect"
	"testing"
)

func TestCache(t *testing.T) {
	// Clean up cache after test
	defer func() {
		mapperCache = newCache()
	}()

	dstTyp := reflect.TypeOf(&[]User{})
	columns := []string{"ID", "Name"}

	m, err := newMapper(dstTyp)
	if err != nil {
		t.Fatalf("error creating new mapper: %s", err)
	}

	mapperCache.storeMap(columns, dstTyp, m)

	loadedMapper, ok := mapperCache.loadMap(columns, dstTyp)
	if !ok {
		t.Fatalf("expected to load mapper from cache, but it was not found")
	}

	if loadedMapper != m {
		t.Errorf("loaded mapper is not the same as the stored mapper")
	}
}
