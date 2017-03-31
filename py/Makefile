SCHEMA=../spec/v2/api.json
SCRIPT=_internal_bin/generate.pl

.PHONY: build

build: $(SCHEMA) $(SCRIPT)
	perl $(SCRIPT) $(SCHEMA) >& octav.py