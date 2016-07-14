SCHEMA=../spec/v1/api.json
SCRIPT=_internal_bin/generate.pl

.PHONY: build

build: $(SCHEMA) $(SCRIPT)
	perl $(SCRIPT) $(SCHEMA) >& octav.py