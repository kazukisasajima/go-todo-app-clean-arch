package entity

type Task struct {
	ID          int 	`json:"id" gorm:"primaryKey"`
	Title       string  `json:"title" gorm:"not null"`
	UserID      int 	`json:"user_id" gorm:"not null"`
}

// jsonの設定を追加しないと、フロントエンドにデータを返す際にKeyがIDなど大文字になってしまう
// 例）
// {
// 	"ID": 1,
// 	"Title": "test",
// 	"UserID": 1
//}

// jsonの設定を追加すると、以下のようになる
// {
// 	"id": 1,
// 	"title": "test",
// 	"user_id": 1
// }