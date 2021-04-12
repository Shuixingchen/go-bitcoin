package bitcoin

type MTree struct {
	Hash string
	Transactions []*Transaction
}

//把当前要打包的交易生成一个merkletree
func CreateMTree(transactions []*Transaction) *MTree{

	return &MTree{
		Hash:         GetSHA256HashCode(transactions),
		Transactions: transactions,
	}
}