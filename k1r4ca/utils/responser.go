package ulits

import "HDUhelper_Todo/models"

func Response(Item []*models.ListItem) []*models.ResponseItem {
	//用ResponseItem结构体框架下的Items进行时间转换 以传输unix时间戳
	Items := make([]*models.ResponseItem, len(Item))
	for i := range Item {
		Items[i] = &models.ResponseItem{
			Id:        Item[i].Id,
			DueDate:   Item[i].DueDate,
			Item:      Item[i].Item,
			Done:      Item[i].Done,
			Over:      Item[i].Over,
			CreatedAt: Item[i].CreatedAt.Unix(),
			UpdatedAt: Item[i].UpdatedAt.Unix(),
		}
	}
	return Items
}
