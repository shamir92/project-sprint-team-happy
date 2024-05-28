# ProjectSprint! Social Media Test Cases

### Prerequisites
- [K6](https://k6.io/docs/get-started/installation/)
- A linux environment (WSL / MacOS should be fine)

### Environment Variables
- `BASE_URL` fill this with your backend url (eg: `http://localhost:8080`)

### How to run
#### For regular testing
```bash
BASE_URL=http://localhost:8080 make run
```
#### For load testing
Prepare the NIP Generator:
```bash
# this requires at minimal, go 1.22.3
# run in separate terminal
make run-go
```

Then run:
```bash
BASE_URL=http://localhost:8080 make run-test-case
```