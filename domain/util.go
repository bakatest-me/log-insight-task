package domain

type KV[k, v any] struct {
	Key   k
	Value v
}
