package serializer

import "whxxxxxxxxxx/model"

type Task struct {
	ID        uint   `json:"id" form:"id" example:"1"`
	UID       uint   `json:"uid" form:"uid" example:"1"`
	Title     string `json:"title" form:"title" example:"任务标题"`
	Content   string `json:"content" form:"content" example:"任务内容"`
	Status    int    `json:"status" form:"status" example:"0"`
	StartTime int64  `json:"start_time" form:"start_time" example:"0"`
	EndTime   int64  `json:"end_time" form:"end_time" example:"0"`
	Username  string `json:"username" form:"username" example:"whxxxxxxxxxx"`
}

func BuildTask(task model.Task) Task {
	return Task{
		ID:        task.ID,
		UID:       task.Uid,
		Username:  task.User.UserName,
		Title:     task.Title,
		Content:   task.Content,
		Status:    task.Status,
		StartTime: task.StartTime,
		EndTime:   task.EndTime,
	}
}

func BuildTasks(items []model.Task) (tasks []Task) {
	for _, item := range items {
		task := BuildTask(item)
		tasks = append(tasks, task)
	}
	return tasks
}
