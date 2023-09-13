
clean:
	rm -f dist/*

stat: clean
	go run -v .
