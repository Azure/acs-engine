@echo off

setlocal enabledelayedexpansion
set DIR=%~dp0
(cd %DIR%.. && go build)
if !ERRORLEVEL! neq 0 goto error
for %%i in (%DIR%*.ank) do (
  %DIR%..\anko %DIR%lib\tester.ank %%i
  if !ERRORLEVEL! neq 0 goto error
)
exit /b 0
:error
exit /b 1
