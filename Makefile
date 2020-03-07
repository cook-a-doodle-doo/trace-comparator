DOTS      := $(wildcard *.dot)
DOTS_NAME := $(patsubst %.dot,%,$(wildcard *.dot))

run:
	go run main.go
	@for obj in $(DOTS_NAME); do \
		dot -T pdf $$obj.dot -o $$obj.pdf; \
	done

#.dot.pdf: 
#	dot -T pdf $< -o $@
