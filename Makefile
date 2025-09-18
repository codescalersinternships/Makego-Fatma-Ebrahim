run: dep1 dep2
	echo "Hello World \n " > hello.txt

clean:
	rm hello
	rm dep1.txt
	rm dep2.txt
dep1: dep2
	echo "Hello World \n " > dep1.txt
dep2: clean
	echo "Hello World \n " > dep2.txt
run:
	echo "no \n " > hello