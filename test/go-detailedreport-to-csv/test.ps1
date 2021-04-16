Remove-Item "$PSScriptRoot\temp" -Recurse -ErrorAction Ignore
Copy-Item -Path "$PSScriptRoot\input" -Destination "$PSScriptRoot\temp" -Recurse
& "$PSScriptRoot\..\..\vcgopkg.exe" "$PSScriptRoot\temp" "-20060102150405"

if( $LASTEXITCODE -eq 0 ) {
	Write-Output "vcgopkg exited successfully"
} else {
    Write-Output "vcgopkg failed"
}

Move-Item -Path $PSScriptRoot\temp -Destination $PSScriptRoot\test
$fsTemp  = Get-ChildItem -Recurse -path $PSScriptRoot\test | Select-Object FullName
Move-Item -Path $PSScriptRoot\test -Destination $PSScriptRoot\temp

Move-Item -Path $PSScriptRoot\output -Destination $PSScriptRoot\test
$fsOutput = Get-ChildItem -Recurse -path $PSScriptRoot\test | Select-Object FullName
Move-Item -Path $PSScriptRoot\test -Destination $PSScriptRoot\output

$c = Compare-Object -ReferenceObject $fsOutput -DifferenceObject $fsOutput

if ($c -ne 0) {
    Write-Output "Compare -ne 0. Compare results:"
    Write-Output $c
    exit 1
} else {
    exit 0
}