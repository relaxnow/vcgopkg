Remove-Item "$PSScriptRoot\temp" -Recurse -ErrorAction Ignore
Copy-Item -Path "$PSScriptRoot\input" -Destination "$PSScriptRoot\temp" -Recurse
& "$PSScriptRoot\..\..\vcgopkg.exe" $PSScriptRoot\temp 20060102150405

Move-Item -Path $PSScriptRoot\temp -Destination $PSScriptRoot\test
$fsTemp  = Get-ChildItem -Recurse -path $PSScriptRoot\test
Move-Item -Path $PSScriptRoot\test -Destination $PSScriptRoot\temp

Move-Item -Path $PSScriptRoot\output -Destination $PSScriptRoot\test
$fsOutput = Get-ChildItem -Recurse -path $PSScriptRoot\test
Move-Item -Path $PSScriptRoot\test -Destination $PSScriptRoot\output

$c = Compare-Object -ReferenceObject $fsOutput -DifferenceObject $fsTemp

if ($c -ne 0) {
    Write-Output $c
    exit 1
}