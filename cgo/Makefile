# Copyright 2013-2023, MACHBASE Corporation or its subsidiaries.
# All rights reserved.

# $Id:
#

include $(MACHBASE_HOME)/install/machbase_env.mk

GCCVERSIONGTEQ7 := $(shell expr `gcc -dumpversion | cut -f1 -d.` \>= 7)

ifeq "$(GCCVERSIONGTEQ7)" "1"
    LD_FLAGS += -no-pie
endif

INCLUDES += $(INC_OPT)/$(MACHBASE_HOME)/include

all : data.txt make_data select append

data.txt : make_data
	./make_data

make_data : make_data.o
	$(LD_CC) $(LD_FLAGS) $(LD_OUT_OPT)$@ $< $(LIB_OPT)machbasecli$(LIB_AFT)  $(LIBDIR_OPT)$(MACHBASE_HOME)/lib $(LD_LIBS)

make_data.o : make_data.c
	$(COMPILE.cc) $(CC_FLAGS) $(INCLUDES) $(CC_OUT_OPT)$@ $<

select : select.go
	go build select.go

append : append.go
	go build append.go

clean :
	rm -f *.o make_data data.txt select append test_result.dat
