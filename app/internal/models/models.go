package models

type Board struct {
	Id string `json:"id" validate:"is_int"`
}

type BoardForAdd struct {
	Title string `json:"title" validate:"min=1,max=150"`
}

type BoardForJoin struct {
	Code string `json:"code"`
}

type User struct {
	Uid string `json:"uid" validate:"min=1,max=128"`
}

type Item struct {
	Id string `validate:"is_int"`
	Y  string `json:"y"`
	W  string `json:"w"`
	T  string `json:"t"`
}

type UserForUpdate struct {
	Uid   string `json:"uid" validate:"min=1,max=128"`
	Name  string `json:"name" validate:"min=1,max=150"`
	Email string `json:"email" validate:"min=1,max=256"`
}
