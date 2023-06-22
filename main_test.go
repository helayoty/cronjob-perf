package main

import (
	"cronjob-perf/pkg/cronjob"
	"fmt"
	"testing"
	"time"

	"github.com/robfig/cron"
)

func BenchmarkNewMostRecentScheduleTime(b *testing.B) {
	sched, _ := cron.ParseStandard("@every 5ms")

	// fmt.Printf("Running the new implementation with schedule %s\n", "*/5 * * * *")  // Runs every 5 minutes
	// cronJobObj1 := cronjob.CreateCronJob("test-cronjob1","*/5 * * * *")
	// for i := 0; i < b.N; i++ {
	// 	cronjob.NewMostRecentScheduleTime(cronJobObj1,time.Now(), sched, false)
	// }

	// fmt.Printf("Running the new implementation with schedule %s\n", "30 6-16/4 * * 1-5")
	// cronJobObj2 := cronjob.CreateCronJob("test-cronjob2", "30 6-16/4 * * 1-5")
	// for i := 0; i < b.N; i++ {
	// 	cronjob.NewMostRecentScheduleTime(cronJobObj2,time.Now(), sched, false)
	// }

	fmt.Printf("Running the new implementation with schedule %s\n", "30 10,11,12 * * 1-5")
	cronJobObj3 := cronjob.CreateCronJob("test-cronjob3", "30 10,11,12 * * 1-5")
	for i := 0; i < b.N; i++ {
		cronjob.NewMostRecentScheduleTime(cronJobObj3, time.Now(), sched, false)
	}
}

func BenchmarkOldMostRecentScheduleTime(b *testing.B) {
	sched, _ := cron.ParseStandard("@every 5ms")

	// fmt.Printf("Running the old implementation with schedule %s\n", "*/5 * * * *")
	// cronJobObj1 := cronjob.CreateCronJob("test-cronjob1","*/5 * * * *")
	// for i := 0; i < b.N; i++ {
	// 	cronjob.OldMostRecentScheduleTime(cronJobObj1,time.Now(), sched, false)
	// }

	// fmt.Printf("Running the old implementation with schedule %s\n", "30 6-16/4 * * 1-5")
	// cronJobObj2 := cronjob.CreateCronJob("test-cronjob2", "30 6-16/4 * * 1-5")
	// for i := 0; i < b.N; i++ {
	// 	cronjob.OldMostRecentScheduleTime(cronJobObj2,time.Now(), sched, false)
	// }

	fmt.Printf("Running the old implementation with schedule %s\n", "30 10,11,12 * * 1-5")
	cronJobObj3 := cronjob.CreateCronJob("test-cronjob3", "30 10,11,12 * * 1-5")
	for i := 0; i < b.N; i++ {
		cronjob.OldMostRecentScheduleTime(cronJobObj3, time.Now(), sched, false)
	}
}
