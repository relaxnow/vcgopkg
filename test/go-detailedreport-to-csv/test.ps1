$fsInput = Get-ChildItem -Recurse -path input
$fsOutput = Get-ChildItem -Recurse -path output
Compare-Object -ReferenceObject $fsInput -DifferenceObject $fsOutput