default: run
run: dep1 dep2
	echo "Hello World \n " > hello.txt

clean: 
	echo "clean" > clean.txt
dep1: clean
	echo "Hello dep1 \n " > dep1.txt
dep2: dep1 clean
	echo "Hello dep2 \n " > dep2.txt
run: dep1 
	echo "no \n " > hello