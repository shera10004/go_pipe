@echo off

rem ====================================

rem --- 입력 파라미터가 없을때 아래 package_name값을 이름으로 할건지 여부 
set dopack=true

rem --- 패키지명(exe) 현재 폴더의 폴더명(패키지명)을 입력하여야 에러가 발생하지 않습니다.
set package_name=middleweb

rem ====================================

set param=%1
set targetdir=%cd%
set godir=%GOPATH%
set option=-v

set do_run=false

if "%param%"=="" (
set do_run=0 
) else (
set do_run=1
set pn=%param%
)

if %do_run%==0 (
if %dopack%==true (
set do_run=2 
set pn=%package_name%
)
)

if not %do_run%==0 (

echo go install %option% ./...
go install %option% ./...

echo move /y %godir%\bin\%pn%.exe %targetdir%
move /y %godir%\bin\%package_name%.exe %targetdir%

echo --- complete

) else (

echo --- 패키지명을 입력하세요 [ 마지막 폴더 이름 : 패키지명]

)





pause