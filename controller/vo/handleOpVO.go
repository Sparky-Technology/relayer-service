package vo

type HandleOpVO struct {
	UserOp          *UserOperationVO `json:"user_op"`
	BeneficiaryAddr string           `json:"beneficiary_addr"`
}

type UserOperationVO struct {
	Sender               string `json:"sender"`
	Nonce                string `json:"nonce"`
	InitCode             string `json:"init_code"`
	CallData             string `json:"call_data"`
	CallGasLimit         string `json:"call_gas_limit"`
	VerificationGasLimit string `json:"verification_gas_limit"`
	PreVerificationGas   string `json:"pre_verification_gas"`
	MaxFeePerGas         string `json:"max_fee_per_gas"`
	MaxPriorityFeePerGas string `json:"max_priority_fee_per_gas"`
	PaymasterAndData     string `json:"paymaster_and_data"`
	Signature            string `json:"signature"`
}
