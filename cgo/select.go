package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>
#include <string.h>
#include <machbase_sqlcli.h>
#include <sys/time.h>

#define MACHBASE_PORT_NO    5656

#define RC_SUCCESS          0
#define RC_FAILURE          -1

#define CHECK_STMT_RESULT(aRC, aSTMT, aMsg)     \
    if( sRC != SQL_SUCCESS )                    \
    {                                           \
        printError(gEnv, gCon, aSTMT, aMsg);    \
        goto error;                             \
    }


SQLHENV     gEnv;
SQLHDBC     gCon;

void printError(SQLHENV aEnv, SQLHDBC aCon, SQLHSTMT aStmt, char *aMsg);
void printColumn(char *aCol, int aLen, char *aFormat, ...);
int connectDB();
void disconnectDB();
int executeDirectSQL(const char *aSQL, int aErrIgnore);
int prepareExecuteSQL(const char *aSQL);
int createTable();
int selectTable();
time_t getTimeStamp();


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

void printColumn(char *aCol, int aLen, char *aFormat, ...)
{
    fprintf(stdout, "%s : ", aCol);

    if( aLen == SQL_NULL_DATA )
    {
        fprintf(stdout, "NULL");
    }
    else
    {
        va_list ap;
        va_start(ap, aFormat);
        vfprintf(stdout, aFormat, ap);
        va_end(ap);
    }
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

int prepareExecuteSQL(const char *aSQL)
{
    SQLHSTMT sStmt = SQL_NULL_HSTMT;

    if( SQLAllocStmt(gCon, &sStmt) != SQL_SUCCESS )
    {
        printError(gEnv, gCon, sStmt, "SQLAllocStmt Error");
        goto error;
    }

    if( SQLPrepare(sStmt, (SQLCHAR *)aSQL, SQL_NTS) != SQL_SUCCESS )
    {
        printError(gEnv, gCon, sStmt, "SQLPrepare Error");
        goto error;
    }

    if( SQLExecute(sStmt) != SQL_SUCCESS )
    {
        printError(gEnv, gCon, sStmt, "SQLExecute Error");
        goto error;
    }

    if( SQLFreeStmt(sStmt, SQL_DROP) != SQL_SUCCESS )
    {
        printError(gEnv, gCon, sStmt, "SQLFreeStmt Error");
        goto error;
    }
    sStmt = SQL_NULL_HSTMT;

    return RC_SUCCESS;

error:
    if( sStmt != SQL_NULL_HSTMT )
    {
        SQLFreeStmt(sStmt, SQL_DROP);
        sStmt = SQL_NULL_HSTMT;
    }

    return RC_FAILURE;
}

int createTable()
{
    int sRC;

    sRC = executeDirectSQL("DROP TABLE CLI_SAMPLE", 1);
    if( sRC != RC_SUCCESS )
    {
        return RC_FAILURE;
    }

    sRC = executeDirectSQL("CREATE TABLE CLI_SAMPLE(seq short, score integer, total long, percentage float, ratio double, id varchar(10), srcip ipv4, dstip ipv6, reg_date datetime, textlog text, image binary)", 0);
    if( sRC != RC_SUCCESS )
    {
        return RC_FAILURE;
    }

    return RC_SUCCESS;
}

int selectTable()
{
    const char *sSQL = "SELECT seq, score, total, percentage, ratio, id, srcip, dstip, reg_date, textlog, image  FROM CLI_SAMPLE";

    SQLHSTMT    sStmt = SQL_NULL_HSTMT;
    SQLRETURN   sRC   = SQL_ERROR;
    //int         i     = 0;

    SQLLEN      sSeqLen     = 0;
    SQLLEN      sScoreLen   = 0;
    SQLLEN      sTotalLen   = 0;
    SQLLEN      sPctLen     = 0;
    SQLLEN      sRatioLen   = 0;
    SQLLEN      sIdLen      = 0;
    SQLLEN      sSrcIPLen   = 0;
    SQLLEN      sDstIPLen   = 0;
    SQLLEN      sRegDateLen = 0;
    SQLLEN      sLogLen     = 0;
    SQLLEN      sImageLen   = 0;

    short                sSeq;
    int                  sScore;
    long                 sTotal;
    float                sPercentage;
    double               sRatio;
    char                 sId[14];
    char                 sSrcIP[16];
    char                 sDstIP[40];
    SQL_TIMESTAMP_STRUCT sRegDate;
    char                 sLog[1024];
    char                 sImage[1024];


    if( SQLAllocStmt(gCon, &sStmt) != SQL_SUCCESS )
    {
        printError(gEnv, gCon, sStmt, "SQLAllocStmt Error");
        goto error;
    }

    if( SQLPrepare(sStmt, (SQLCHAR *)sSQL, SQL_NTS) != SQL_SUCCESS )
    {
        printError(gEnv, gCon, sStmt, "SQLPrepare Error");
        goto error;
    }

    if( SQLExecute(sStmt) != SQL_SUCCESS )
    {
        printError(gEnv, gCon, sStmt, "SQLExecute Error");
        goto error;
    }

    sRC = SQLBindCol(sStmt, 1, SQL_C_SHORT, &sSeq, 0, &sSeqLen);
    CHECK_STMT_RESULT(sRC, sStmt, "SQLBindCol 1 Error");

    sRC = SQLBindCol(sStmt, 2, SQL_C_LONG, &sScore, 0, &sScoreLen);
    CHECK_STMT_RESULT(sRC, sStmt, "SQLBindCol 2 Error");

    sRC = SQLBindCol(sStmt, 3, SQL_C_BIGINT, &sTotal, 0, &sTotalLen);
    CHECK_STMT_RESULT(sRC, sStmt, "SQLBindCol 3 Error");

    sRC = SQLBindCol(sStmt, 4, SQL_C_FLOAT, &sPercentage, 0, &sPctLen);
    CHECK_STMT_RESULT(sRC, sStmt, "SQLBindCol 4 Error");

    sRC = SQLBindCol(sStmt, 5, SQL_C_DOUBLE, &sRatio, 0, &sRatioLen);
    CHECK_STMT_RESULT(sRC, sStmt, "SQLBindCol 5 Error");

    sRC = SQLBindCol(sStmt, 6, SQL_C_CHAR, sId, sizeof(sId), &sIdLen);
    CHECK_STMT_RESULT(sRC, sStmt, "SQLBindCol 6 Error");

    sRC = SQLBindCol(sStmt, 7, SQL_C_CHAR, sSrcIP, sizeof(sSrcIP), &sSrcIPLen);
    CHECK_STMT_RESULT(sRC, sStmt, "SQLBindCol 7 Error");

    sRC = SQLBindCol(sStmt, 8, SQL_C_CHAR, sDstIP, sizeof(sDstIP), &sDstIPLen);
    CHECK_STMT_RESULT(sRC, sStmt, "SQLBindCol 8 Error");

    sRC = SQLBindCol(sStmt, 9, SQL_C_TYPE_TIMESTAMP, &sRegDate, 0, &sRegDateLen);
    CHECK_STMT_RESULT(sRC, sStmt, "SQLBindCol 9 Error");

    sRC = SQLBindCol(sStmt, 10, SQL_C_CHAR, sLog, sizeof(sLog), &sLogLen);
    CHECK_STMT_RESULT(sRC, sStmt, "SQLBindCol 10 Error");

    sRC = SQLBindCol(sStmt, 11, SQL_C_CHAR, sImage, sizeof(sImage), &sImageLen);
    CHECK_STMT_RESULT(sRC, sStmt, "SQLBindCol 11 Error");

    while( SQLFetch(sStmt) == SQL_SUCCESS ){}

    if( SQLFreeStmt(sStmt, SQL_DROP) != SQL_SUCCESS )
    {
        printError(gEnv, gCon, sStmt, "SQLFreeStmt Error");
        goto error;
    }
    sStmt = SQL_NULL_HSTMT;

    return RC_SUCCESS;

error:
    if( sStmt != SQL_NULL_HSTMT )
    {
        SQLFreeStmt(sStmt, SQL_DROP);
        sStmt = SQL_NULL_HSTMT;
    }

    return RC_FAILURE;
}

time_t getTimeStamp()
{
    struct timeval sTimeVal;

    gettimeofday(&sTimeVal, NULL);

    return (sTimeVal.tv_sec*10000000 + sTimeVal.tv_usec);
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
	"fmt"
)

func main() {

	if C.connectDB() == C.RC_SUCCESS {
		fmt.Println("connectDB success.")
	} else {
		fmt.Println("connectDB failure.")
		appendError()
	}

	var sStartTime, sEndTime C.time_t
	sStartTime = C.getTimeStamp()
	if C.selectTable() != C.RC_SUCCESS {
		fmt.Println("selectTable failure.")
		appendError()
	}
	sEndTime = C.getTimeStamp()

	fmt.Println("selectTable success")
	fmt.Println(fmt.Sprintf("%.2f", ((C.double)(sEndTime-sStartTime))/10000000), "second")

	C.disconnectDB()

	C.fflush(C.stdout)

}

func appendError() {

	if C.checkCon() != C.RC_SUCCESS {
		C.disconnectDB()
	}

}
