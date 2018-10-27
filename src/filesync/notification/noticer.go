package notification

import (
	"fmt"
)

func Email(object string,  subject string,mailTo []string, noticeContent string) (bool, string) {
	fmt.Println(object, subject, mailTo, noticeContent)
	return true, ""
}

func Sms(object string, smsList []string, noticeContent string ) (bool, string){
	fmt.Println(object, smsList, noticeContent)
	return true, ""
}


