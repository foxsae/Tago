taggo: taggo.6
	6l -o taggo taggo.6
taggo.6: 
	6g taggo.go
clean:
	rm taggo taggo.6
.PHONY: clean
