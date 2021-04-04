Copy-Item -Path "$PSScriptRoot\input" -Destination "$PSScriptRoot\temp" -Recurse
& "$PSScriptRoot\..\..\vcgopkg.exe" $PSScriptRoot\temp
$fsInput = Get-ChildItem -Recurse -path $PSScriptRoot\temp
$fsOutput = Get-ChildItem -Recurse -path $PSScriptRoot\output
Compare-Object -ReferenceObject $fsInput -DifferenceObject $fsOutput