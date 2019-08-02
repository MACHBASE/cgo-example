#!/bin/bash

rm test_result.dat
rm data.txt
./make_data

COUNT=10
AVERAGE=0.0

for ((i=1;i<=$COUNT;i++)); do
	echo "======================= append $i =======================" | tee -a test_result.dat
	APPEND_RES=`./append | tail -n 3 | grep -v "appendClose success"`
	echo "$APPEND_RES" | tee -a test_result.dat
	SECOND=`echo "$APPEND_RES" | grep second | awk '{print $1}'`
	AVERAGE=`echo "$AVERAGE $SECOND" | awk '{print $1+$2}'`
	echo "" | tee -a test_result.dat
done

AVERAGE=`echo "$AVERAGE $COUNT" | awk '{printf "%.2f", $1/$2}'`
echo "append average : $AVERAGE" | tee -a test_result.dat
echo "" | tee -a test_result.dat
echo "" | tee -a test_result.dat

AVERAGE=0.0

for ((i=1;i<=$COUNT;i++)); do
	echo "======================= select $i =======================" | tee -a test_result.dat
	SELECT_RES=`./select | tail -n 1`
	echo "$SELECT_RES" | tee -a test_result.dat
	SECOND=`echo "$SELECT_RES" | grep second | awk '{print $1}'`
	AVERAGE=`echo "$AVERAGE $SECOND" | awk '{print $1+$2}'`
	echo "" | tee -a test_result.dat
done

AVERAGE=`echo "$AVERAGE $COUNT" | awk '{printf "%.2f", $1/$2}'`
echo "" | tee -a test_result.dat
echo "select average : $AVERAGE" | tee -a test_result.dat
