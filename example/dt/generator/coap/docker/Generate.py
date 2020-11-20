#! /usr/bin/env python[3.4.0]
# -*- coding: utf-8 -*-
import time
from pytz import timezone
from datetime import datetime
import random



# 신한 크레인 데이터 포맷으로 구성된 문자열 데이터 입니다.
data_list = [
            b"H 135.2000.0-07.6000.0364.9000.001090101999.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9000.000000000.00000000000000000010000000000",
                b"H 135.2000.0-07.6000.0364.9000.001030101999.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9000.000000000.00000000000000000010000000000",
                    b"H 135.2000.0-07.6000.0364.9000.001090101999.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9000.000000000.00000000000000000010000000000",
                        b"H 135.2000.0-07.6000.0364.9000.001090101999.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9000.000000000.00000000000000000010000000000"
                        ]


# 데이터 생성 함수입니다.
def RandomGenerate(CraneFullName):
    # 현재 시간 구하기
    ctime = datetime.now(timezone('Asia/Seoul')).strftime("%Y-%m-%d %H:%M:%S")

    # 데이터 포맷은 업체_크레인명, 버전, 현재시간, 데이터 문자열로 구성됩니다.
    DataString = CraneFullName + '|'
    DataString += ctime+"|"

    # 신한 크레인 포맷의 stx, etx 추가
    stx = b"\x02"
    etx = b"\x03"

    # 랜덤하게 데이터 리스트로부터 데이터를 얻어서 신한 크레인 포맷 데이터 구성
    data = data_list[random.randint(0,4)]

    DataString += data.decode()
    print("["+ctime + "] '" + CraneFullName + "' Crane Random Data Generate")
    print("--> "+ DataString)
    return DataString
