# Week 4 Requirements and K6 Load Test Results

## Requirements

### This Week's Tasks

1. Please open Belimang!.mhtml on browser to view this week's requirements.
2. K6 run on 0.25vcpu and 0.5GB RAM AWS ECS Fargate for Golang Apps and Postgre DB.

## Postman Collection

1. Download the postman collection for this project [here](./ProjectSprint%20-%20BeliMang.postman_collection.json).
2. Open Postman.
3. Click on the **Import** button located at the top left of the Postman application.
4. Select the `Import File` tab.
5. Click on **Choose Files** and select the provided Postman collection file (`Project Sprint Batch 2.postman_collection.json`).
6. Click **Import**.

## K6 Load Test Results

### Summary

- **Total Checks:** 10717
- **Successful Checks:** 10715 (99.98%)
- **Failed Checks:** 2 (0.02%)

### Data Transfer

- **Data Received:** 36 MB (421 kB/s)
- **Data Sent:** 2.6 MB (30 kB/s)

### gRPC Request Duration

- **Average:** 5.1 ms
- **Minimum:** 234.26 µs
- **Median:** 409.07 µs
- **Maximum:** 1.06 s
- **90th Percentile:** 750.15 µs
- **95th Percentile:** 4.96 ms

### HTTP Request Metrics

- **Blocked:**

  - Average: 57.63 µs
  - Minimum: 1.25 µs
  - Median: 4.08 µs
  - Maximum: 49.02 ms
  - 90th Percentile: 6.71 µs
  - 95th Percentile: 398.09 µs

- **Connecting:**

  - Average: 46.76 µs
  - Minimum: 0 µs
  - Median: 0 µs
  - Maximum: 48.94 ms
  - 90th Percentile: 0 µs
  - 95th Percentile: 313.14 µs

- **Duration:**

  - Average: 2.67 s
  - Minimum: 1.15 ms
  - Median: 717 ms
  - Maximum: 37.2 s
  - 90th Percentile: 8.13 s
  - 95th Percentile: 13.4 s

- **Failed Requests:** 8.20% (399 out of 4862)

- **Receiving:**

  - Average: 86.49 µs
  - Minimum: 21.33 µs
  - Median: 55.44 µs
  - Maximum: 58.83 ms
  - 90th Percentile: 100.32 µs
  - 95th Percentile: 116.75 µs

- **Sending:**

  - Average: 85.27 µs
  - Minimum: 8.5 µs
  - Median: 26.85 µs
  - Maximum: 50.67 ms
  - 90th Percentile: 55 µs
  - 95th Percentile: 71.98 µs

- **TLS Handshaking:** 0 µs (no TLS handshakes)

- **Waiting:**
  - Average: 2.67 s
  - Minimum: 1.05 ms
  - Median: 716.13 ms
  - Maximum: 37.2 s
  - 90th Percentile: 8.13 s
  - 95th Percentile: 13.4 s

### General Metrics

- **Total HTTP Requests:** 4862 (57.006575/s)
- **Iteration Duration:**

  - Average: 15.04 s
  - Minimum: 2.44 ms
  - Median: 7.27 s
  - Maximum: 1m 5s
  - 90th Percentile: 38.34 s
  - 95th Percentile: 52.01 s

- **Virtual Users:**

  - Minimum: 0
  - Maximum: 300

- **Maximum Virtual Users:** 300

## Notes

- The above metrics indicate the performance and reliability of the system under load.
- Address any failed checks and optimize accordingly to improve overall performance.
- Ensure all requirements are met by the end of the week in preparation for the final project.

---

For any questions or assistance, please reach out to the project manager.
