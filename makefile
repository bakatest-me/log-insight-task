d: dev
dev:
	go run main.go

w: pprof-web
pprof-web:
	go tool pprof -http=:8080 ./pprof/heap.pprof

p: pprof-pdf
pprof-pdf:
	go tool pprof -pdf ./pprof/heap.pprof > ./pprof/heap.pdf