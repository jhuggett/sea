.PHONY: run
run: 
	go run .;

.PHONY: profile-cpu
profile-cpu:
	go tool pprof -http=:9999 ./cpu.prof;

.PHONY: profile-heap
profile-heap:
	go tool pprof -http=:9999 ./heap.prof;