@echo off
chcp 65001 >nul
for /f %%A in ('powershell -NoProfile -Command ""') do rem sync console
setlocal enabledelayedexpansion

:: --- 定位到脚本所在目录 ---
cd /d %~dp0

:: --- 自动跳转到 go.mod 根目录 ---
for /f "usebackq tokens=*" %%i in (`go env GOMOD 2^>NUL`) do set GOMOD=%%i
if "%GOMOD%"=="" (
  echo [ERR] go.mod not found. Run in module root. go mod init xxx
  exit /b 1
)
for %%i in ("%GOMOD%") do set MODROOT=%%~dpi
cd /d "%MODROOT%"

:: --- 你的配置 ---
set APP_NAME=jd-worklog
set MAIN_PKG=./cmd/server
set CGO_ENABLED=0
set CC_386=i686-linux-gnu-gcc
set CC_AMD64=x86_64-linux-gnu-gcc
set EXTRA_BUILD_TAGS=

:: --- 版本信息 ---
set VERSION=%~1
if "%VERSION%"=="" set VERSION=dev
for /f "tokens=* usebackq" %%i in (`git rev-parse --short=8 HEAD 2^>NUL`) do set GIT_COMMIT=%%i
if "%GIT_COMMIT%"=="" set GIT_COMMIT=nogit
for /f "tokens=* usebackq" %%i in (`git describe --tags --abbrev=0 2^>NUL`) do set GIT_TAG=%%i
if "%GIT_TAG%"=="" set GIT_TAG=none
for /f "tokens=* usebackq" %%i in (`powershell -NoProfile -Command "(Get-Date).ToUniversalTime().ToString('yyyy-MM-ddTHH:mm:ssZ')"`) do set BUILD_TIME=%%i
set LDFLAGS=-s -w -X main.Version=%VERSION% -X main.GitCommit=%GIT_COMMIT% -X main.GitTag=%GIT_TAG% -X main.BuildTime=%BUILD_TIME%

:: --- 预检：MAIN_PKG 是否存在源码 ---
if not exist "%MAIN_PKG%\*.go" (
  echo [ERR] No .go files under "%CD%\%MAIN_PKG%".
  echo   set MAIN_PKG=./cmd/server
  exit /b 1
)

:: --- 输出目录 ---
set DIST=dist
if not exist "%DIST%" mkdir "%DIST%"

echo ----------------------------------------
echo Building %APP_NAME% version %VERSION%
echo Git: commit=%GIT_COMMIT% tag=%GIT_TAG%
echo Time: %BUILD_TIME%
echo CGO_ENABLED=%CGO_ENABLED%
echo Tags=%EXTRA_BUILD_TAGS%
echo OutDir=%DIST%
echo ----------------------------------------

:: ===== linux/386 =====
set GOOS=linux
set GOARCH=386
if "%CGO_ENABLED%"=="1" set CC=%CC_386%
call :build_one linux_386 || exit /b 1

:: ===== linux/amd64 =====
set GOOS=linux
set GOARCH=amd64
if "%CGO_ENABLED%"=="1" set CC=%CC_AMD64%
call :build_one linux_amd64 || exit /b 1

powershell -NoProfile -Command "Compress-Archive -Path '%DIST%\*' -DestinationPath '%DIST%\%APP_NAME%_%VERSION%.zip' -Force" 1>NUL 2>NUL
echo echo [Pack] Creating ZIP archive...: %DIST%
exit /b 0

:build_one
set SUFFIX=%~1
set OUT=%DIST%\%APP_NAME%_%SUFFIX%
if "%EXTRA_BUILD_TAGS%"=="" (set TAGS_OPT=) else (set TAGS_OPT=-tags "%EXTRA_BUILD_TAGS%")
echo GOOS=%GOOS% GOARCH=%GOARCH% CGO_ENABLED=%CGO_ENABLED% CC=%CC%
go build -trimpath -ldflags "%LDFLAGS%" %TAGS_OPT% -o "%OUT%" "%MAIN_PKG%" || (
  echo [ERR] Build failed: %SUFFIX%
  exit /b 1
)
certutil -hashfile "%OUT%" SHA256 > "%OUT%.sha256.txt" 2>NUL
echo [OK] Build success: %OUT%
exit /b 0
