package builder

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/nervosnetwork/ckb-sdk-go/systemscript"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

type sudtMockIterator struct {
	Cells []*types.TransactionInput
	index int
}

func (m *sudtMockIterator) HasNext() bool {
	return m.index < len(m.Cells)
}

func (m *sudtMockIterator) Next() *types.TransactionInput {
	current := m.Cells[m.index]
	m.index += 1
	return current
}

var (
	sudtSender, _ = address.Decode("ckt1qzda0cr08m85hc8jlnfp3zer7xulejywt49kt2rr0vthywaa50xwsq02cgdvd5mng9924xarf3rflqzafzmzlpsuhh83c")
	sudtArgs, _   = hexutil.Decode("0xae4147ba8412767b3fd9bd16d45dab2fa5df283a6fd68dae5367524daa767ca7")
	sudtType      = &types.Script{
		CodeHash: types.HexToHash("0xc5e5dcf215925f7ef4dfaf5f4b4f105bc321c02776d6e7d52a1db3fcd9d011a4"),
		HashType: types.HashTypeType,
		Args:     sudtArgs,
	}
)

func getSudtMockIterator() *sudtMockIterator {
	return &sudtMockIterator{
		Cells: []*types.TransactionInput{
			{
				OutPoint: &types.OutPoint{
					TxHash: types.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
					Index:  0,
				},
				Output: &types.CellOutput{
					Capacity: 100000000000,
					Lock:     sudtSender.Script,
					Type:     sudtType,
				},
				OutputData: systemscript.EncodeSudtAmount(big.NewInt(100)),
			},
			{
				OutPoint: &types.OutPoint{
					TxHash: types.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
					Index:  1,
				},
				Output: &types.CellOutput{
					Capacity: 10000000000,
					Lock:     sudtSender.Script,
					Type:     sudtType,
				},
				OutputData: systemscript.EncodeSudtAmount(big.NewInt(10)),
			},
		},
	}
}

func TestSudtTransactionBuilderBalance(t *testing.T) {
	var err error
	iterator := getSudtMockIterator()
	builder := NewSudtTransactionBuilderFromSudtArgs(types.NetworkTest, iterator, SudtTransactionTypeTransfer, sudtArgs)

	if _, err = builder.AddSudtOutputByAddress("ckt1qzda0cr08m85hc8jlnfp3zer7xulejywt49kt2rr0vthywaa50xwsqdamwzrffgc54ef48493nfd2sd0h4cjnxg4850up", big.NewInt(1)); err != nil {
		t.Error(err)
	}
	builder.FeeRate = 1000
	if err = builder.AddChangeOutputByAddress("ckt1qzda0cr08m85hc8jlnfp3zer7xulejywt49kt2rr0vthywaa50xwsqdamwzrffgc54ef48493nfd2sd0h4cjnxg4850up"); err != nil {
		t.Error(err)
	}
	tx, err := builder.Build()
	if err != nil {
		t.Error(err)
	}

	amount1, err := systemscript.DecodeSudtAmount(tx.TxView.OutputsData[0])
	if err != nil {
		t.Error(err)
	}
	amount2, err := systemscript.DecodeSudtAmount(tx.TxView.OutputsData[1])
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, big.NewInt(100), amount1.Add(amount1, amount2))
}
