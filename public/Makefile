
exec:object
	ld -o $(FILE) $(FILE).o

lib:object
	ld -dynamic-linker /lib/ld-linux.so.2 -o $(FILE) -lc $(FILE).o

object:$(FILE).s
	as -gstabs -o $(FILE).o $(FILE).s
.PHONY:clean
clean:
	-rm $(FILE) $(FILE).o
