@echo off

rem ====================================

rem --- �Է� �Ķ���Ͱ� ������ �Ʒ� package_name���� �̸����� �Ұ��� ���� 
set dopack=true

rem --- ��Ű����(exe) ���� ������ ������(��Ű����)�� �Է��Ͽ��� ������ �߻����� �ʽ��ϴ�.
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

echo --- ��Ű������ �Է��ϼ��� [ ������ ���� �̸� : ��Ű����]

)





pause