package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/jinzhu/now"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var calendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "print calendar events",
	Run:   calendar,
}

var (
	dbPath      string
	startDate   string
	withRestDay bool
)

func init() {
	homeDir, _ := os.UserHomeDir()
	calendarCmd.Flags().StringVarP(&dbPath, "db", "", fmt.Sprintf("%s/Library/Calendars/Calendar.sqlitedb", homeDir), "SQLite数据库文件路径")
	calendarCmd.Flags().StringVarP(&startDate, "start", "s", now.BeginningOfMonth().Format("2006-01-02"), "开始时间")
	calendarCmd.Flags().BoolVarP(&withRestDay, "rest", "r", false, "统计休息日加班")

	rootCmd.AddCommand(calendarCmd)
}

type CalendarItem struct {
	Summary   string
	StartDate string
	EndTDate  string
}

// NSTimeIntervalSince1970 978307200 = 2001-01-01 00:00:00 UTC
const NSTimeIntervalSince1970 = 978307200

func calendar(cmd *cobra.Command, args []string) {
	switch runtime.GOOS {
	case "darwin":
	default:
		log.Fatal("只支持 macOS 系统")
	}

	db, err := gorm.Open(sqlite.Open(dbPath))
	if err != nil {
		cobra.CheckErr(err)
	}

	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		cobra.CheckErr(err)
	}
	startCalTime := startTime.Add(-NSTimeIntervalSince1970 * time.Second)

	keyword := "summary like '加班%小时'"
	if withRestDay {
		keyword = "summary like '加班%天'"
	}

	var items []CalendarItem
	err = db.Table("CalendarItem").
		Select(fmt.Sprintf("summary, %s,  %s", toUnixTimeColumn("start_date"), toUnixTimeColumn("end_date"))).
		Where(keyword).
		Where("start_date >= ?", startCalTime.Unix()).
		Order("start_date").
		Find(&items).Error
	if err != nil {
		cobra.CheckErr(err)
	}

	if len(items) == 0 {
		log.Print("没有找到加班记录")
		os.Exit(0)
	}

	var data [][]string
	count := 0
	for _, v := range items {
		t, _ := time.Parse("2006-01-02 15:04:05", v.StartDate)

		data = append(data, []string{
			fmtDate(t),
			fmtSummary(v)},
		)

		count += cast.ToInt(strings.Trim(v.Summary, "加班小时"))
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"日期", "时长"})
	table.SetFooter([]string{fmt.Sprintf("%d天", len(items)), fmt.Sprintf("%d小时", count)}) // Add Footer
	table.SetBorder(false)                                                                // Set Border to false
	table.AppendBulk(data)                                                                // Add Bulk Data
	table.Render()                                                                        // Send output
}

func toUnixTimeColumn(columnName string) string {
	return fmt.Sprintf("datetime(%s + %d, 'unixepoch', 'localtime') as %[1]s", columnName, NSTimeIntervalSince1970)
}

func fmtSummary(v CalendarItem) string {
	return strings.Trim(v.Summary, "加班")
}

func fmtDate(t time.Time) string {
	return fmt.Sprintf("%s %s",
		t.Format("01-02"),
		weekdayToChinese(t.Weekday()))
}

// weekday to chinese
func weekdayToChinese(weekday time.Weekday) string {
	switch weekday {
	case time.Sunday:
		return "日"
	case time.Monday:
		return "一"
	case time.Tuesday:
		return "二"
	case time.Wednesday:
		return "三"
	case time.Thursday:
		return "四"
	case time.Friday:
		return "五"
	case time.Saturday:
		return "六"
	default:
		return ""
	}
}
