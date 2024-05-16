package utils

import "time"

func GetDayTime(from int) (startTime, endTime int64) {
	switch from {
	case 0: // 今天
		dateNow := time.Now()
		startTime = time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 0, 0, 0, 0, dateNow.Location()).UnixNano() / 1e6
		endTime = time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 23, 59, 59, 0, dateNow.Location()).UnixNano() / 1e6
	//case "1": // 昨天
	//	dateNow := time.Now().AddDate(0, 0, -1)
	//	startTime = time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 0, 0, 0, 0, dateNow.Location()).UnixNano() / 1e6
	//	endTime = time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 23, 59, 59, 0, dateNow.Location()).UnixNano() / 1e6
	//case "2": // 前天
	//	dateNow := time.Now().AddDate(0, 0, -2)
	//	startTime = time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 0, 0, 0, 0, dateNow.Location()).UnixNano() / 1e6
	//	endTime = time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 23, 59, 59, 0, dateNow.Location()).UnixNano() / 1e6
	//case "3":
	//	dateNow := time.Now().AddDate(0, 0, -3)
	//	startTime = time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 0, 0, 0, 0, dateNow.Location()).UnixNano() / 1e6
	//	endTime = time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 23, 59, 59, 0, dateNow.Location()).UnixNano() / 1e6
	//case "4":
	//	dateNow := time.Now().AddDate(0, 0, -4)
	//	startTime = time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 0, 0, 0, 0, dateNow.Location()).UnixNano() / 1e6
	//	endTime = time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 23, 59, 59, 0, dateNow.Location()).UnixNano() / 1e6
	//case "5":
	//	dateNow := time.Now().AddDate(0, 0, -5)
	//	startTime = time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 0, 0, 0, 0, dateNow.Location()).UnixNano() / 1e6
	//	endTime = time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 23, 59, 59, 0, dateNow.Location()).UnixNano() / 1e6
	//case "6":
	//	dateNow := time.Now().AddDate(0, 0, -6)
	//	startTime = time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 0, 0, 0, 0, dateNow.Location()).UnixNano() / 1e6
	//	endTime = time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 23, 59, 59, 0, dateNow.Location()).UnixNano() / 1e6
	default:
		dateNow := time.Now().AddDate(0, 0, -from)
		startTime = time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 0, 0, 0, 0, dateNow.Location()).UnixNano() / 1e6
		endTime = time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 23, 59, 59, 0, dateNow.Location()).UnixNano() / 1e6
	}
	return
}

// from: 1: 1d, 2: 2d
func GetFromNowTime(from int) (startTime, endTime int64) {
	endTime = time.Now().UnixNano() / 1e6
	startTime = time.Now().AddDate(0, 0, -from).UnixNano() / 1e6
	return
}
