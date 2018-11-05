cls
go install .\Model
go install .\View
go install .\Controller

go build -ldflags="-H windowsgui -s -w"

@set NewFileName=".\Peer001\Distributed ledger demo.exe"
@xcopy /F /Y .\1_DistributedLedgerDemo.exe %NewFileName%
@REM @set NewFileName=".\Peer002\Distributed ledger demo.exe"
@REM @xcopy /F /Y .\1_DistributedLedgerDemo.exe %NewFileName%
@set NewFileName="Distributed ledger demo.exe"
@move .\1_DistributedLedgerDemo.exe %NewFileName%
