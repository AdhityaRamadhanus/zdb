package zdb

func Map[In, Out any](inputs []In, fn func(In) Out) []Out {
	res := []Out{}
	for _, input := range inputs {
		res = append(res, fn(input))
	}
	return res
}

func Reduce[In any, Acc any](inputs []In, fn func(Acc, In) Acc, init Acc) (acc Acc) {
	acc = init
	for _, input := range inputs {
		acc = fn(acc, input)
	}
	return acc
}
