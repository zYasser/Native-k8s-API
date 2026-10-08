@echo off
setlocal
cd /d "%~dp0"

rem Apply only the .yml/.yaml files in the root directory (no subdirectories)
for %%f in (*.yml *.yaml) do (
    echo Applying %%f
    kubectl apply -f "%%f"
)

endlocal
