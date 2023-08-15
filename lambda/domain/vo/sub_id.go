package vo

type SubId string

func NewSubId(v string) SubId {
	return SubId(v)
}
