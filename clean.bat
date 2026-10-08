@echo off
setlocal
cd /d "%~dp0"

rem Delete every resource defined by a .yml/.yaml file in this project (recursively), ignoring vendor and .git
for /r %%f in (*.yml *.yaml) do call :process "%%f"

endlocal
goto :eof

:process
set "p=%~1"
if not "%p:\vendor\=%"=="%p%" goto :eof
if not "%p:\.git\=%"=="%p%" goto :eof
echo Deleting resources from %p%
kubectl delete -f "%p%" --ignore-not-found
goto :eof
