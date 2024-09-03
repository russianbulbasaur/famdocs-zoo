build:
	mkdir outputs
	go build -o outputs/server cmd/famdocs/main.go
	mkdir outputs/uploads
	cp .env outputs/.env

clean:
	rm -rf uploads/*
	rm -rf outputs/*
