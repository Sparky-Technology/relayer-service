# relayer-service

## Introduction

**Relayer Service** is backend service for swapbotAA

## Project setup
```
go build main.go
```

### Config
set config file 'config.yaml' under conf directory
```
ethereum:
  rpc_endpoint_addr: 
  wallet_contract_addr: 
  uniswap_router_contract_addr: 
  system_account_private_key: 
```

### Start
start relayer service
```
./main
```
