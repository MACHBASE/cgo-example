package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <arpa/inet.h>
#include <sys/time.h>
#include <machbase_sqlcli.h>

#define MACHBASE_PORT_NO	5656
#define ERROR_CHECK_COUNT	100

#define RC_SUCCESS          0
#define RC_FAILURE          -1

#define UNUSED(aVar) do { (void)(aVar); } while(0)

#define CHECK_STMT_RESULT(aRC, aSTMT, aMsg)     \
    if( sRC != SQL_SUCCESS )                    \
    {                                           \
        printError(gEnv, gCon, aSTMT, aMsg);    \
        goto error;                             \
    }


SQLHENV 	gEnv;
SQLHDBC 	gCon;

void printError(SQLHENV aEnv, SQLHDBC aCon, SQLHSTMT aStmt, char *aMsg);
time_t getTimeStamp();
int connectDB();
void disconnectDB();
int executeDirectSQL(const char *aSQL, int aErrIgnore);
int createTable();
int appendOpen(SQLHSTMT aStmt);
int appendData(SQLHSTMT aStmt, char* seq, char* score, char* total, char* percentage, char* ratio, char* id, char* srcip, char* dstip, char* datetime, char* textlog, char* binary);
SQLBIGINT appendClose(SQLHSTMT aStmt);

void printError(SQLHENV aEnv, SQLHDBC aCon, SQLHSTMT aStmt, char *aMsg)
{
    SQLINTEGER      sNativeError;
    SQLCHAR         sErrorMsg[SQL_MAX_MESSAGE_LENGTH + 1];
    SQLCHAR         sSqlState[SQL_SQLSTATE_SIZE + 1];
    SQLSMALLINT     sMsgLength;

    if( aMsg != NULL )
    {
        printf("%s\n", aMsg);
    }

    if( SQLError(aEnv, aCon, aStmt, sSqlState, &sNativeError,
        sErrorMsg, SQL_MAX_MESSAGE_LENGTH, &sMsgLength) == SQL_SUCCESS )
    {
        printf("SQLSTATE-[%s], Machbase-[%d][%s]\n", sSqlState, sNativeError, sErrorMsg);
    }
}

void appendDumpError(SQLHSTMT    aStmt,
                 SQLINTEGER  aErrorCode,
                 SQLPOINTER  aErrorMessage,
                 SQLLEN      aErrorBufLen,
                 SQLPOINTER  aRowBuf,
                 SQLLEN      aRowBufLen)
{
    char       sErrMsg[1024] = {0, };
    char       sRowMsg[32 * 1024] = {0, };

    UNUSED(aStmt);

    if (aErrorMessage != NULL)
    {
        strncpy(sErrMsg, (char *)aErrorMessage, aErrorBufLen);
    }

    if (aRowBuf != NULL)
    {
        strncpy(sRowMsg, (char *)aRowBuf, aRowBufLen);
    }

    fprintf(stdout, "Append Error : [%d][%s]\n[%s]\n\n", aErrorCode, sErrMsg, sRowMsg);
}

time_t getTimeStamp()
{
    struct timeval sTimeVal;

    gettimeofday(&sTimeVal, NULL);

    return (sTimeVal.tv_sec*10000000 + sTimeVal.tv_usec);
}


int connectDB()
{
    char sConnStr[1024];

    if( SQLAllocEnv(&gEnv) != SQL_SUCCESS )
    {
        printf("SQLAllocEnv error\n");
        return RC_FAILURE;
    }

    if( SQLAllocConnect(gEnv, &gCon) != SQL_SUCCESS )
    {
        printf("SQLAllocConnect error\n");

        SQLFreeEnv(gEnv);
        gEnv = SQL_NULL_HENV;

        return RC_FAILURE;
    }

    sprintf(sConnStr,"SERVER=127.0.0.1;UID=SYS;PWD=MANAGER;CONNTYPE=1;PORT_NO=%d", MACHBASE_PORT_NO);

    if( SQLDriverConnect( gCon, NULL,
                          (SQLCHAR *)sConnStr,
                          SQL_NTS,
                          NULL, 0, NULL,
                          SQL_DRIVER_NOPROMPT ) != SQL_SUCCESS
      )
    {

        printError(gEnv, gCon, NULL, "SQLDriverConnect error");

        SQLFreeConnect(gCon);
        gCon = SQL_NULL_HDBC;

        SQLFreeEnv(gEnv);
        gEnv = SQL_NULL_HENV;

        return RC_FAILURE;
    }

    return RC_SUCCESS;
}

void disconnectDB()
{
    if( SQLDisconnect(gCon) != SQL_SUCCESS )
    {
        printError(gEnv, gCon, NULL, "SQLDisconnect error");
    }

    SQLFreeConnect(gCon);
    gCon = SQL_NULL_HDBC;

    SQLFreeEnv(gEnv);
    gEnv = SQL_NULL_HENV;
}

int executeDirectSQL(const char *aSQL, int aErrIgnore)
{
    SQLHSTMT sStmt = SQL_NULL_HSTMT;

    if( SQLAllocStmt(gCon, &sStmt) != SQL_SUCCESS )
    {
        if( aErrIgnore == 0 )
        {
            printError(gEnv, gCon, sStmt, "SQLAllocStmt Error");
            return RC_FAILURE;
        }
    }

    if( SQLExecDirect(sStmt, (SQLCHAR *)aSQL, SQL_NTS) != SQL_SUCCESS )
    {

        if( aErrIgnore == 0 )
        {
            printError(gEnv, gCon, sStmt, "SQLExecDirect Error");

            SQLFreeStmt(sStmt,SQL_DROP);
            sStmt = SQL_NULL_HSTMT;
            return RC_FAILURE;
        }
    }

    if( SQLFreeStmt(sStmt, SQL_DROP) != SQL_SUCCESS )
    {
        if (aErrIgnore == 0)
        {
            printError(gEnv, gCon, sStmt, "SQLFreeStmt Error");
            sStmt = SQL_NULL_HSTMT;
            return RC_FAILURE;
        }
    }
    sStmt = SQL_NULL_HSTMT;

    return RC_SUCCESS;
}

int createTable()
{
    int sRC;

    sRC = executeDirectSQL("DROP TABLE CLI_SAMPLE", 1);
    if( sRC != RC_SUCCESS )
    {
        return RC_FAILURE;
    }

    sRC = executeDirectSQL("CREATE TABLE CLI_SAMPLE(seq short, score integer, total long, percentage float, ratio double, id varchar(13), srcip ipv4, dstip ipv6, reg_date datetime, textlog text, image binary)", 0);
    if( sRC != RC_SUCCESS )
    {
        return RC_FAILURE;
    }

    return RC_SUCCESS;
}

int appendOpen(SQLHSTMT aStmt)
{
    const char *sTableName = "CLI_SAMPLE";

    if( SQLAppendOpen(aStmt, (SQLCHAR *)sTableName, ERROR_CHECK_COUNT) != SQL_SUCCESS )
    {
        printError(gEnv, gCon, aStmt, "SQLAppendOpen Error");
        return RC_FAILURE;
    }

    return RC_SUCCESS;
}


SQLBIGINT        sCount;
int appendData(SQLHSTMT aStmt, char* seq, char* score, char* total, char* percentage, char* ratio, char* id, char* srcip, char* dstip, char* datetime, char* textlog, char* binary)
{
    SQL_APPEND_PARAM sParam[11];
    memset(sParam, 0, sizeof(sParam));
    sParam[0].mShort = atoi(seq);   //short
    sParam[1].mInteger = atoi(score); //int
    sParam[2].mLong = atoi(total);    //long
    sParam[3].mFloat = atoi(percentage);   //float
    sParam[4].mDouble = atoi(ratio);  //double
    sParam[5].mVar.mLength = strlen(id);
    sParam[5].mVar.mData = id;
    sParam[6].mIP.mLength = SQL_APPEND_IP_STRING;
    sParam[6].mIP.mAddrString = srcip;
    sParam[7].mIP.mLength = SQL_APPEND_IP_STRING;
    sParam[7].mIP.mAddrString = dstip;
    sParam[8].mDateTime.mTime = SQL_APPEND_DATETIME_STRING;
    sParam[8].mDateTime.mDateStr = datetime;
    sParam[8].mDateTime.mFormatStr = "DD/MON/YYYY:HH24:MI:SS";
    sParam[9].mVar.mLength = strlen(textlog);
    sParam[9].mVar.mData = textlog;
    sParam[10].mVar.mLength = strlen(binary);
    sParam[10].mVar.mData = binary;
    if( SQLAppendDataV2(aStmt, sParam) != SQL_SUCCESS )
    {
            printError(gEnv, gCon, aStmt, "SQLAppendData Error");
            return RC_FAILURE;
    }
    if ( ((sCount++) % 10000) == 0)
    {
            fflush(stdout);
    }
    if( ((sCount) % 100) == 0 )
    {
            if( SQLAppendFlush( aStmt ) != SQL_SUCCESS )
            {
                    printError(gEnv, gCon, aStmt, "SQLAppendFlush Error");
            }
    }
    return RC_SUCCESS;
}

SQLBIGINT appendClose(SQLHSTMT aStmt)
{
    SQLBIGINT sSuccessCount = 0;
    SQLBIGINT  sFailureCount = 0;

    if( SQLAppendClose(aStmt, &sSuccessCount, &sFailureCount) != SQL_SUCCESS )
    {
        printError(gEnv, gCon, aStmt, "SQLAppendClose Error");
        return RC_FAILURE;
    }

    printf("success : %ld, failure : %ld\n", sSuccessCount, sFailureCount);

    return sSuccessCount;
}

void initStmt(SQLHSTMT* sStmt) {
	*sStmt = SQL_NULL_HSTMT;
}

int checkStmt(SQLHSTMT* sStmt) {
	if(*sStmt != SQL_NULL_HSTMT) {
		return RC_FAILURE;
	}
	return RC_SUCCESS;
}

int checkCon() {
	if(gCon != SQL_NULL_HDBC) {
		return RC_FAILURE;
	}
	return RC_SUCCESS;
}
*/
import "C"
import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	var sStmt C.SQLHSTMT
	var sStartTime C.time_t
	var sEndTime C.time_t

	C.initStmt(&sStmt)

	if C.connectDB() == C.RC_SUCCESS {
		fmt.Println("connectDB success")
	} else {
		appendError(&sStmt)
	}
	if C.createTable() == C.RC_SUCCESS {
		fmt.Println("createTable success")
	} else {
		fmt.Println("createTable failure")
		appendError(&sStmt)
	}
	if C.SQLAllocStmt(C.gCon, &sStmt) != C.SQL_SUCCESS {
		C.printError(C.gEnv, C.gCon, sStmt, C.CString("SQLAllocStmt Error"))
		appendError(&sStmt)
	}
	if C.appendOpen(sStmt) == C.RC_SUCCESS {
		fmt.Println("appendOpen success")
	} else {
		fmt.Println("appendOpen failure")
		appendError(&sStmt)
	}

	data, _ := os.Open("data.txt")
	defer data.Close()
	scanner := bufio.NewScanner(data)

	sStartTime = C.getTimeStamp()
	C.sCount = 0
	for scanner.Scan() {
		s := scanner.Text()
		row := strings.Split(s, ",")
		C.appendData(sStmt, C.CString(row[0]), C.CString(row[1]), C.CString(row[2]), C.CString(row[3]), C.CString(row[4]), C.CString(row[5]), C.CString(row[6]), C.CString(row[7]), C.CString(row[8]), C.CString(row[9]), C.CString(row[10]))
	}
	sEndTime = C.getTimeStamp()

	C.sCount = C.appendClose(sStmt)
	if C.sCount >= 0 {
		C.fflush(C.stdout)
		fmt.Println("appendClose success")
		fmt.Println(fmt.Sprintf("%.2f", ((C.double)(sEndTime-sStartTime))/10000000), "second")
	} else {
		fmt.Println("appendClose failure")
	}
	if C.SQLFreeStmt(sStmt, C.SQL_DROP) != C.SQL_SUCCESS {
		C.printError(C.gEnv, C.gCon, sStmt, C.CString("SQLFreeStmt Error"))
		appendError(&sStmt)
	}
	C.initStmt(&sStmt)

	C.disconnectDB()

}

func appendError(sStmt *C.SQLHSTMT) {

	if C.checkStmt(sStmt) != C.RC_SUCCESS {
		C.SQLFreeStmt(*sStmt, C.SQL_DROP)
		C.initStmt(sStmt)
	}

	if C.checkCon() != C.RC_SUCCESS {
		C.disconnectDB()
	}

}
