package custmodels

type ListCommon struct {
	Page  uint64 `json:"-"`
	Limit uint64 `json:"-"`
}