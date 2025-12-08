run:
	go tool templ generate --watch --proxy="http://localhost:3000" --cmd="go run ."
