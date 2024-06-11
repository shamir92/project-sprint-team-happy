# BeliMang! ProjectSprint Batch 2 Week 4 Test cases!

## Prerequisites
- [ k6 ](https://k6.io/docs/get-started/installation/)
- golang (min ver 1.22.3)

## Guide
### Run locally
1. Open two terminal
2. Run the caddy backend in one of the terminal by using:
    ```bash
    make run
    ```
3. And if you ready to test, run:
    ```bash
    make test
    ```
#### Debugging
If you want to know why your test fail, you can use:
```bash
make test-debug
```

> 💡 `make test` & `make test-debug` will automatically set the BASE_URL of the backend to `http://localhost:8080`, if you wish to change it, you can use:
>```bash
>BASE_URL=http://backend.url:8080 make test
>```
