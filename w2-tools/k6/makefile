
.PHONY: run_debug
run_debug:
	git pull origin main && DEBUG_ALL=true k6 run script.js > output.txt 2>&1

.PHONY: run_timed
run_timed:
	git pull origin main && k6 run --duration 10s script.js

.PHONY: run
run:
	git pull origin main && k6 run script.js

.PHONY: runLoadTest
runLoadTest:
	git pull origin main && LOAD_TEST=true k6 run script.js