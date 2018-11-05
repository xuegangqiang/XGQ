cls
@go install .\KademliaInformationSystem
@go build -ldflags="-s -w"
@xcopy /F /Y .\2_KademliaSimu.exe .\Svr\KademliaSimu.exe
@move .\2_KademliaSimu.exe .\Clt\KademliaSimu.exe
