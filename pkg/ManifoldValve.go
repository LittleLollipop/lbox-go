package lbox

type ManifoldValveInterface interface {
	OnValveOpen(data map[string]interface{}) error
}

type ManifoldValve struct {
	valveList []string
	outputImp ManifoldValveInterface
	data      map[string]interface{}
}

func NewManifoldValve(valveList []string, outputImp ManifoldValveInterface) *ManifoldValve {
	return &ManifoldValve{
		valveList: valveList,
		outputImp: outputImp,
		data:      make(map[string]interface{}),
	}
}

func (mv *ManifoldValve) Input(valveName string, data interface{}) {
	for _, iterator := range mv.valveList {
		if iterator == valveName {
			mv.data[valveName] = data
		}
	}
	mv.CheckValve()
}

func (mv *ManifoldValve) CheckValve() {
	readyKeys := make([]string, 0, len(mv.data))
	ready := 0
	for key := range mv.data {
		readyKeys = append(readyKeys, key)
	}
	for _, iterator := range mv.valveList {
		if contains(readyKeys, iterator) {
			ready++
		}
	}
	if len(mv.valveList) == ready {
		mv.outputImp.OnValveOpen(mv.data)
	}
}

func contains(keys []string, key string) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}
