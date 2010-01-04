tao: tao.6
	6l -o tao tao.6
tao.6: 
	6g tao.go
clean:
	rm tao tao.6
.PHONY: clean
