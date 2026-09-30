package main

import "fmt"

type Appointment struct{ ID, PatientName, Phone, StartsAt, Status string }

func reminderText(a Appointment) string {
	return fmt.Sprintf("Hi %s, your appointment %s starts at %s. Reply to your clinic if you need help.", a.PatientName, a.ID, a.StartsAt)
}

func SendReminder(client *SMSClient, a Appointment) (string, error) {
	if a.Status != "confirmed" {
		return "", nil
	}
	return client.Send(a.Phone, reminderText(a), "appointment-reminder-"+a.ID)
}
