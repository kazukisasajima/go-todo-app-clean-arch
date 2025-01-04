package entity

func NewDomains() []interface{} {
	return []interface{}{&Task{}, &User{}}
}