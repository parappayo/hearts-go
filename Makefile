
run:
	go run cmd/cli/main.go

serve:
	go run cmd/server/main.go

test:
	go test hearts/pkg/cards
	go test hearts/pkg/game
