package enums

const (
	App        string = "go-hexagonal-scheduler"
	PostfixDev string = "dev"

	TaskNameOne     string = "TaskEveryMinute"
	TaskScheduleOne string = "@every 1m"
	TaskRunEveryOne int    = 60000

	TaskNameTwo     string = "TaskEveryTwoMinutes"
	TaskScheduleTwo string = "@every 2m"
	TaskRunEveryTwo int    = 120000
)
