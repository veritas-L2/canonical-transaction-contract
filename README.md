# canonical-transaction-contract


## Start the chaincode

```
CORE_CHAINCODE_LOGLEVEL=debug CORE_PEER_TLS_ENABLED=false CORE_CHAINCODE_ID_NAME=ctc:1.0 ./ctc -peer.address 127.0.0.1:7052
```

## Approve and commit the chaincode definition

```
peer lifecycle chaincode approveformyorg  -o 127.0.0.1:7050 --channelID ch1 --name ctc --version 1.0 --sequence 1 --init-required --signature-policy "OR ('SampleOrg.member')" --package-id ctc:1.0
peer lifecycle chaincode checkcommitreadiness -o 127.0.0.1:7050 --channelID ch1 --name ctc --version 1.0 --sequence 1 --init-required --signature-policy "OR ('SampleOrg.member')"
peer lifecycle chaincode commit -o 127.0.0.1:7050 --channelID ch1 --name ctc --version 1.0 --sequence 1 --init-required --signature-policy "OR ('SampleOrg.member')" --peerAddresses 127.0.0.1:7051
```

## Using the contract 

```
CORE_PEER_ADDRESS=127.0.0.1:7051 peer chaincode invoke -o 127.0.0.1:7050 -C ch1 -n ctc -c '{"Args":["Init"]}' --isInit

CORE_PEER_ADDRESS=127.0.0.1:7051 peer chaincode invoke -o 127.0.0.1:7050 -C ch1 -n ctc -c '{"Args":["GetBatchHistory"]}'

CORE_PEER_ADDRESS=127.0.0.1:7051 peer chaincode invoke -o 127.0.0.1:7050 -C ch1 -n ctc -c '{"Args":["CommitBatch", "{\"timestamp\":1670195170,\"transactions\":[{\"chaincodeName\":\"basic\",\"transactionName\":\"ReadAsset\",\"args\":[\"goo\"]},{\"chaincodeName\":\"basic\",\"transactionName\":\"ReadAsset\",\"args\":[\"goo\"]},{\"chaincodeName\":\"basic\",\"transactionName\":\"ReadAsset\",\"args\":[\"goo\"]},{\"chaincodeName\":\"basic\",\"transactionName\":\"ReadAsset\",\"args\":[\"goo\"]},{\"chaincodeName\":\"basic\",\"transactionName\":\"ReadAsset\",\"args\":[\"goo\"]},{\"chaincodeName\":\"basic\",\"transactionName\":\"ReadAsset\",\"args\":[\"goo\"]},{\"chaincodeName\":\"basic\",\"transactionName\":\"ReadAsset\",\"args\":[\"goo\"]},{\"chaincodeName\":\"basic\",\"transactionName\":\"ReadAsset\",\"args\":[\"goo\"]},{\"chaincodeName\":\"basic\",\"transactionName\":\"ReadAsset\",\"args\":[\"goo\"]},{\"chaincodeName\":\"basic\",\"transactionName\":\"ReadAsset\",\"args\":[\"goo\"]}],\"prevStateHash\":\"c27+1x/yZplDJ6BTyP5qRxO2N9ndqMvHTfMaiV46prI=\",\"newStateHash\":\"c27+1x/yZplDJ6BTyP5qRxO2N9ndqMvHTfMaiV46prI=\"}"]}'

CORE_PEER_ADDRESS=127.0.0.1:7051 peer chaincode invoke -o 127.0.0.1:7050 -C ch1 -n ctc -c '{"Args":["DeleteBatchChain"]}'
```