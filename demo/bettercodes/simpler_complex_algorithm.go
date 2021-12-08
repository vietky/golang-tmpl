package bettercodes

type SimplerComplexAlgorithm struct {
	ComplexAlgorithm
}

func (sc *SimplerComplexAlgorithm) DoComplexThings() int {
	return sc.ComplexAlgorithm.complex()
}

func NewSimplerComplexAlgorithm() IComplexAlgorithm {
	return &SimplerComplexAlgorithm{}
}
