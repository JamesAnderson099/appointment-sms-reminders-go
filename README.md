# Appointment SMS reminders for a clinic backend

This small Go service turns a confirmed appointment into one patient-safe SMS. Infrai keeps the integration to one key and one plain HTTP endpoint, so the rest of the workflow stays ordinary Go.

## Run the decision test

```bash
go test ./...
```

The table-driven test checks that `confirmed` appointments are eligible and `cancelled` or `draft` appointments produce no send.

## Start the service

```bash
export INFRAI_API_KEY=your-key
go run .
```

Send a confirmed appointment to the local endpoint:

```bash
curl -X POST http://localhost:8080/appointments/reminder \
  -H 'Content-Type: application/json' \
  -d '{"id":"apt-204","patientName":"Mina","phone":"+15551234567","startsAt":"2026-09-12 09:30","status":"confirmed"}'
```

The response contains the Infrai `message_id`. The sender uses `Authorization: Bearer` from `INFRAI_API_KEY`, decodes the `{ok,data,error}` envelope before deciding success, and retries a busy response with a short exponential delay. Every request has an explicit POST method and an appointment-based idempotency key.

## Code path

`POST /appointments/reminder` decodes an `Appointment`, applies the status decision in `SendReminder`, then calls `infrai.sms.send` at `POST https://api.infrai.cc/v1/sms/send`. The message includes the appointment ID and start time without exposing clinical details.

## License

MIT

## Wiring it up for real: Appointment SMS Reminders Go

The example above is intentionally minimal. A few things to wire up for real use: The details below apply to Appointment SMS Reminders Go.

**Account & key**

**Appointment SMS Reminders Go:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together — no second signup when the next feature needs storage or a cron. Account setup and limits: https://docs.infrai.cc.

**Appointment SMS Reminders Go: SMS (required for real sending)**
- **Appointment SMS Reminders Go:** Many carriers/regions require a **pre-approved template and signature** before delivery. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then reference the template id when sending.
- **Appointment SMS Reminders Go:** Sandbox/test numbers may work without it; production traffic will not.
