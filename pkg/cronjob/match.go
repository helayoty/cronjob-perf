package cronjob

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	batchv1 "k8s.io/api/batch/v1"
)

func NewMostRecentScheduleTime(cj *batchv1.CronJob, now time.Time, schedule cron.Schedule, includeStartingDeadlineSeconds bool) (time.Time, *time.Time, int64, error) {
	earliestTime := cj.ObjectMeta.CreationTimestamp.Time
	if cj.Status.LastScheduleTime != nil {
		earliestTime = cj.Status.LastScheduleTime.Time
	}
	if includeStartingDeadlineSeconds && cj.Spec.StartingDeadlineSeconds != nil {
		// controller is not going to schedule anything below this point
		schedulingDeadline := now.Add(-time.Second * time.Duration(*cj.Spec.StartingDeadlineSeconds))
		if schedulingDeadline.After(earliestTime) {
			earliestTime = schedulingDeadline
		}
	}
	t1 := schedule.Next(earliestTime)
	t2 := schedule.Next(t1)
	if now.Before(t1) {
		return earliestTime, nil, 0, nil
	}
	if now.Before(t2) {
		return earliestTime, &t1, 1, nil
	}
	// It is possible for cron.ParseStandard("59 23 31 2 *") to return an invalid schedule
	// minute - 59, hour - 23, dom - 31, month - 2, and dow is optional, clearly 31 is invalid
	// In this case the timeBetweenTwoSchedules will be 0, and we error out the invalid schedule
	timeBetweenTwoSchedules := int64(t2.Sub(t1).Round(time.Second).Seconds())
	if timeBetweenTwoSchedules < 1 {
		return earliestTime, nil, 0, fmt.Errorf("time difference between two schedules is less than 1 second")
	}
	// this logic used for calculating number of missed schedules does a rough
	// approximation, by calculating a diff between two schedules (t1 and t2),
	// and counting how many of these will fit in between last schedule and now
	timeElapsed := int64(now.Sub(t1).Seconds())
	numberOfMissedSchedules := (timeElapsed / timeBetweenTwoSchedules) + 1

	var mostRecentTime time.Time
	// to get the most recent time accurate for regular schedules and the ones
	// specified with @every form we first need to calculate potential earliest
	// time by multipliying number of missed schedules by its interval, this
	// is critical to ensure @every starts at the correct time, this explains
	// the numberOfMissedSchedules-1, the additional -1 serves there to go back
	// in time one more time unit back and let the cron library calculate
	// a proper schedule, for case where the schedule is not consistent,
	// for example something like  30 6-16/4 * * 1-5
	potentialEarliest := t1.Add(time.Duration((numberOfMissedSchedules-1-1)*timeBetweenTwoSchedules) * time.Second)
	for t := schedule.Next(potentialEarliest); !t.After(now); t = schedule.Next(t) {
		mostRecentTime = t
	}

	return earliestTime, &mostRecentTime, numberOfMissedSchedules, nil
}

func OldMostRecentScheduleTime(cj *batchv1.CronJob, now time.Time, schedule cron.Schedule, includeStartingDeadlineSeconds bool) (time.Time, *time.Time, int64, error) {
	earliestTime := cj.ObjectMeta.CreationTimestamp.Time
	if cj.Status.LastScheduleTime != nil {
		earliestTime = cj.Status.LastScheduleTime.Time
	}
	if includeStartingDeadlineSeconds && cj.Spec.StartingDeadlineSeconds != nil {
		// controller is not going to schedule anything below this point
		schedulingDeadline := now.Add(-time.Second * time.Duration(*cj.Spec.StartingDeadlineSeconds))
		if schedulingDeadline.After(earliestTime) {
			earliestTime = schedulingDeadline
		}
	}
	t1 := schedule.Next(earliestTime)
	t2 := schedule.Next(t1)
	if now.Before(t1) {
		return earliestTime, nil, 0, nil
	}
	if now.Before(t2) {
		return earliestTime, &t1, 1, nil
	}
	// It is possible for cron.ParseStandard("59 23 31 2 *") to return an invalid schedule
	// minute - 59, hour - 23, dom - 31, month - 2, and dow is optional, clearly 31 is invalid
	// In this case the timeBetweenTwoSchedules will be 0, and we error out the invalid schedule
	timeBetweenTwoSchedules := int64(t2.Sub(t1).Round(time.Second).Seconds())
	if timeBetweenTwoSchedules < 1 {
		return earliestTime, nil, 0, fmt.Errorf("time difference between two schedules is less than 1 second")
	}

	var numberOfMissedSchedules int64
	var mostRecentTime time.Time
	for t := schedule.Next(earliestTime); !t.After(now); t = schedule.Next(t) {
		numberOfMissedSchedules++
		mostRecentTime = t
	}

	return earliestTime, &mostRecentTime, numberOfMissedSchedules, nil
}
