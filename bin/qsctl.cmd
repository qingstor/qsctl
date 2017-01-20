@echo OFF
REM="""
setlocal
set PythonExe=
set PythonExeFlags=

for %%i in (cmd bat exe) do (
    for %%j in (python.%%i) do (
        call :SetPythonExe "%%~$PATH:j"
    )
)
for /f "tokens=2 delims==" %%i in ('assoc .py') do (
    for /f "tokens=2 delims==" %%j in ('ftype %%i') do (
        for /f "tokens=1" %%k in ("%%j") do (
            call :SetPythonExe %%k
        )
    )
)
"%PythonExe%" -x %PythonExeFlags% "%~f0" %*
goto :EOF

:SetPythonExe
if not [%1]==[""] (
    if ["%PythonExe%"]==[""] (
        set PythonExe=%~1
    )
)
goto :EOF
"""

# ===================================================
# Python script starts here
# ===================================================

#!/usr/bin/env python
# -*- coding: utf-8 -*-
# =========================================================================
# Copyright (C) 2016 Yunify, Inc.
# -------------------------------------------------------------------------
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this work except in compliance with the License.
# You may obtain a copy of the License in the LICENSE file, or at:
#
#  http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# =========================================================================

import sys

import qingstor.qsctl.driver

def main():
    return qingstor.qsctl.driver.main()

if __name__ == '__main__':
    sys.exit(main())
