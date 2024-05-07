
.PHONY: run_debug
run_debug:
	DEBUG_ALL=true k6 run script.js

.PHONY: run_timed
run_timed:
	k6 run --duration 10s script.js

.PHONY: run
run:
	k6 run script.js

.PHONY: run50
run50:
	LOAD_TEST=true k6 run -i 50 --vus 50 script.js | grep -E 'checks|http_req_duration|\{ expected_response:true \}'

.PHONY: run100
run100:
	LOAD_TEST=true k6 run -i 100 --vus 100 script.js | grep -E 'checks|http_req_duration|\{ expected_response:true \}'

.PHONY: run200
run200:
	LOAD_TEST=true k6 run -i 200 --vus 200 script.js | grep -E 'checks|http_req_duration|\{ expected_response:true \}'

.PHONY: run300
run300:
	LOAD_TEST=true k6 run -i 300 --vus 300 script.js | grep -E 'checks|http_req_duration|\{ expected_response:true \}'

.PHONY: run600
run600:
	LOAD_TEST=true k6 run -i 600 --vus 600 script.js | grep -E 'checks|http_req_duration|\{ expected_response:true \}'

.PHONY: runAllLoadTests
runAllLoadTests:
	make run50 && sleep 5 && make run100 && sleep 5 && make run200 && sleep 5 && make run300 && sleep 5 && make run600
