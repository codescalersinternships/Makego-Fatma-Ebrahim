default: run
run: dep1 dep2
	echo "Hello World \n " > hello.txt   
clean: 
	rm -rf hello dep1.txt dep2.txt 
dep1: 
	echo "Hello dep1 \n " > dep1.txt
dep2: dep1 
	echo "Hello dep2 \n " > dep2.txt
run: dep1 
	echo "no \n " > hello   