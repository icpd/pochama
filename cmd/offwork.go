package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	workTime.Flags().StringVarP((*string)(&onWorkTime), "onwork", "o", "09:00", "上班时间")
	rootCmd.AddCommand(workTime)
}

var onWorkTime WorkTime

const dinnerTime = time.Minute * 30

var workTime = &cobra.Command{
	Use:   "offwork",
	Short: "print offwork time",
	Run: func(cmd *cobra.Command, args []string) {
		if onWorkTime == "" {
			log.Fatal("上班时间不能为空")
		}

		stdOffWorkTime := onWorkTime.Time().Add(9 * time.Hour)
		afterOneHourOffWork := stdOffWorkTime.Add(dinnerTime).Add(1 * time.Hour)
		afterTwoHourOffWork := stdOffWorkTime.Add(dinnerTime).Add(2 * time.Hour)
		fmt.Printf("0h: %s\n1h：%s\n2h：%s\n",
			stdOffWorkTime.Format("15:04"),
			afterOneHourOffWork.Format("15:04"),
			afterTwoHourOffWork.Format("15:04"))
	},
}

type WorkTime string

func (w WorkTime) Time() time.Time {
	t, _ := time.Parse("15:04", string(w))
	return t
}
