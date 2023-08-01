SUBDIRS := $(wildcard */.)

$(SUBDIRS):
	$(MAKE) -C $@ $(MAKECMDGOALS)

mocks: $(SUBDIRS)

cleanmocks: $(SUBDIRS)

test: $(SUBDIRS)

cleantest: $(SUBDIRS)

db: $(SUBDIRS)

cleandb: $(SUBDIRS)

migrate-new: 
	make -C ${SERVICE} $@ 

migrate-up: $(SUBDIRS)

migrate-down: $(SUBDIRS)

service: $(SUBDIRS)

cleanservice: $(SUBDIRS)

tidy: $(SUBDIRS)

.PHONY: mocks cleanmocks test cleantest db cleandb flyway service cleanservice $(SUBDIRS)
