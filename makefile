tago: tago.6
	6l -o tago tago.6
tago.6: 
	6g tago.go
clean:
	rm tago tago.6
.PHONY: clean
