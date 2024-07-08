# ProjecetSprint Batch 2

[ProjectSprint](https://openidea-projectsprint.notion.site/ProjectSprint-Batch-2-27641ac94ab34fde9b1b97933f8a01f3) is an intensive, project-based learning program designed to bridge the gap between theoretical knowledge and real-world backend development skills.

## Projects

-   Cats Social

    -   [Functional & Non Functional Requirement](https://openidea-projectsprint.notion.site/Cats-Social-9e7639a6a68748c38c67f81d9ab3c769?pvs=4)

    -   [Github Repository](./w1/)

    -   [Postman Collection](./docs/w1/ProjectSprint%20-%20CatSocial.postman_collection.json)

    -   Load Test Results

        | Metric   | 50 VUs | 100 VUs | 200 VUs | 300 VUs | 600 VUs |
        | -------- | ------ | ------- | ------- | ------- | ------- |
        | Checks   | 99.23% | 99.23%  | 99.23%  | 99.23%  | 99.23%  |
        | Avg (ms) | 253.68 | 491.81  | 843.53  | 1210    | 2320    |
        | Min (ms) | 1.82   | 2.15    | 2.13    | 2.58    | 2.26    |
        | Med (ms) | 14.2   | 37.26   | 49.29   | 84.96   | 142.9   |
        | Max (ms) | 2920   | 5560    | 10910   | 16400   | 34810   |
        | p95 (ms) | 2470   | 4970    | 8660    | 11320   | 20740   |
        | p99 (ms) | 2820   | 5380    | 9450    | 14920   | 30000   |

-   Eniqilo Store

    -   [Functional & Non Functional Requirement](https://openidea-projectsprint.notion.site/EniQilo-Store-93d69f62951c4c8aaf91e6c090127886?pvs=4)

    -   [Github Repository](./w2/)

    -   [Postman Collection](./docs/w2/ProjectSprint%20-%20%20EniQilo%20Store.postman_collection.json)

    -   Load Test Results
        | Metric | Value |
        | ----------------- | -------- |
        | Checks | 82% |
        | http_req_failed | 26.33% |
        | data_received | 21 MB |
        | data_sent | 6.8 MB |
        | avg response time | 827.35ms |
        | min response time | 1.95ms |
        | med response time | 482.74ms |
        | max response time | 24.84s |
        | p90 response time | 1.78s |
        | p95 response time | 2.65s |
    -   Docker Hub: [https://hub.docker.com/repository/docker/enjaytarigan99/eniqilo-team-happy](https://hub.docker.com/repository/docker/enjaytarigan99/eniqilo-team-happy)

-   HaloSuster

    -   [Functional & Non Functional Requirement](https://openidea-projectsprint.notion.site/HaloSuster-be1866776fe84c2d8d9eac08ce09b7a5?pvs=4)

    -   [Github Repository](./w3/)

    -   [Postman Collection](./docs/w3/ProjectSprint%20-%20%20HaloSuster.postman_collection.json)

    -   Load Test Results

        | Metric              | Value    |
        | ------------------- | -------- |
        | Checks              | 92.20%   |
        | Requests per second | 769      |
        | Data received       | 71 MB    |
        | Data sent           | 25 MB    |
        | Avg response time   | 819.27ms |
        | Min response time   | 57.03ms  |
        | Med response time   | 484.75ms |
        | Max response time   | 31930ms  |
        | p90 response time   | 1140ms   |
        | p95 response time   | 1990ms   |

-   BeliMang

    -   [Functional & Non Functional Requirement](https://openidea-projectsprint.notion.site/BeliMang-7979300c7ce54dbf8ecd0088806eff14?pvs=4)

    -   [Github Repository](./w4/)

    -   [Postman Collection](./docs/w4/ProjectSprint%20-%20BeliMang.postman_collection.json)

    -   Load Test Results

        | Metric                 | Value               |
        | ---------------------- | ------------------- |
        | Checks                 | 99.98%              |
        | HTTP Requests Failed   | 8.20%               |
        | Requests per second    | 57.006575           |
        | Data received          | 36 MB (421 kB/s)    |
        | Data sent              | 2.6 MB (30 kB/s)    |
        | Avg HTTP req duration  | 2.67s               |
        | Min HTTP req duration  | 1.15ms              |
        | Med HTTP req duration  | 717ms               |
        | Max HTTP req duration  | 37.2s               |
        | p90 HTTP req duration  | 8.13s               |
        | p95 HTTP req duration  | 13.4s               |
        | Avg iteration duration | 15.04s              |
        | Total iterations       | 5                   |
        | Virtual Users (VUs)    | 52 (min=0, max=300) |
        | Max Virtual Users      | 300                 |

        **Additional Information:**

        -   GRPC request duration (avg): 5.1ms
        -   HTTP req blocked (avg): 57.63µs
        -   HTTP req connecting (avg): 46.76µs
        -   HTTP req receiving (avg): 86.49µs
        -   HTTP req sending (avg): 85.27µs
        -   HTTP req waiting (avg): 2.67s

## Team Members

-   [shamir92](https://github.com/shamir92)
-   [14mdzk](https://github.com/14mdzk)
-   [enjaytarigan](https://github.com/enjaytarigan)
