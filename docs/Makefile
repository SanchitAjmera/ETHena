EXTRA_SRC =
TEX = pdflatex

all: Report.pdf

.phony: all clean

%.pdf: %.tex $(EXTRA_SRC)
	while( $(TEX) $< ; grep -q "Rerun to get cross" $(<:.tex=.log)) ; \
	  do true ; \
	 done

clean:
	rm -f *.aux
	rm -f *.log
	rm -f *.pdf
