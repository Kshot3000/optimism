package genesis

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	gethparams "github.com/ethereum/go-ethereum/params"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"

	"github.com/ethereum-optimism/optimism/op-core/eip1559"
	"github.com/ethereum-optimism/optimism/op-core/params"
	"github.com/ethereum-optimism/optimism/op-core/predeploys"
	"github.com/ethereum-optimism/optimism/op-service/ptr"
)

// forkTimes is an OP fork schedule for a fixture; nil leaves a fork unscheduled.
type forkTimes struct {
	regolith, canyon, ecotone, fjord, granite, holocene, isthmus, jovian, karst, lagoon *uint64
}

var (
	staggeredForks = forkTimes{
		regolith: ptr.New[uint64](0), canyon: ptr.New[uint64](2000), ecotone: ptr.New[uint64](3000),
		fjord: ptr.New[uint64](4000), granite: ptr.New[uint64](5000), holocene: ptr.New[uint64](6000),
		isthmus: ptr.New[uint64](7000), jovian: ptr.New[uint64](8000),
	}
	isthmusAtGenesis = forkTimes{
		regolith: ptr.New[uint64](0), canyon: ptr.New[uint64](0), ecotone: ptr.New[uint64](0),
		fjord: ptr.New[uint64](0), granite: ptr.New[uint64](0), holocene: ptr.New[uint64](0),
		isthmus: ptr.New[uint64](0), jovian: ptr.New[uint64](0),
	}
)

func withKarst(f forkTimes, karst uint64) forkTimes {
	f.karst = &karst
	return f
}

func withLagoon(f forkTimes, lagoon uint64) forkTimes {
	f.lagoon = &lagoon
	return f
}

func opChainConfig(f forkTimes) *params.ChainConfig {
	return &params.ChainConfig{
		ChainID: big.NewInt(901),
		Optimism: &params.OptimismConfig{
			EIP1559Elasticity:        6,
			EIP1559Denominator:       50,
			EIP1559DenominatorCanyon: ptr.New[uint64](250),
		},
		BedrockBlock: big.NewInt(0),
		RegolithTime: f.regolith,
		CanyonTime:   f.canyon,
		EcotoneTime:  f.ecotone,
		FjordTime:    f.fjord,
		GraniteTime:  f.granite,
		HoloceneTime: f.holocene,
		IsthmusTime:  f.isthmus,
		JovianTime:   f.jovian,
		KarstTime:    f.karst,
		LagoonTime:   f.lagoon,
	}
}

// opGethChainConfig is the op-geth chain config that op-chain-ops writes into an L2
// genesis.json for the fork schedule f.
func opGethChainConfig(f forkTimes) *gethparams.ChainConfig {
	return &gethparams.ChainConfig{
		ChainID:                 big.NewInt(901),
		HomesteadBlock:          big.NewInt(0),
		EIP150Block:             big.NewInt(0),
		EIP155Block:             big.NewInt(0),
		EIP158Block:             big.NewInt(0),
		ByzantiumBlock:          big.NewInt(0),
		ConstantinopleBlock:     big.NewInt(0),
		PetersburgBlock:         big.NewInt(0),
		IstanbulBlock:           big.NewInt(0),
		MuirGlacierBlock:        big.NewInt(0),
		BerlinBlock:             big.NewInt(0),
		LondonBlock:             big.NewInt(0),
		ArrowGlacierBlock:       big.NewInt(0),
		GrayGlacierBlock:        big.NewInt(0),
		MergeNetsplitBlock:      big.NewInt(0),
		TerminalTotalDifficulty: big.NewInt(0),
		BedrockBlock:            big.NewInt(0),
		RegolithTime:            f.regolith,
		CanyonTime:              f.canyon,
		ShanghaiTime:            f.canyon,
		CancunTime:              f.ecotone,
		EcotoneTime:             f.ecotone,
		FjordTime:               f.fjord,
		GraniteTime:             f.granite,
		HoloceneTime:            f.holocene,
		IsthmusTime:             f.isthmus,
		JovianTime:              f.jovian,
		KarstTime:               f.karst,
		PragueTime:              f.isthmus,
		LagoonTime:              f.lagoon,
		Optimism: &gethparams.OptimismConfig{
			EIP1559Elasticity:        6,
			EIP1559Denominator:       50,
			EIP1559DenominatorCanyon: ptr.New[uint64](250),
		},
	}
}

func fixtureAlloc() types.GenesisAlloc {
	return types.GenesisAlloc{
		predeploys.L2ToL1MessagePasserAddr: {
			Code:    []byte{0x60, 0x00, 0x60, 0x00, 0xf3},
			Storage: map[common.Hash]common.Hash{{0x01}: {0x02}, {0x03}: {0x04}},
			Balance: big.NewInt(0),
			Nonce:   1,
		},
		common.HexToAddress("0x00000000000000000000000000000000000beef0"): {
			Balance: new(big.Int).Lsh(big.NewInt(1), 200),
		},
	}
}

// fixture returns the same genesis as a Genesis and as the op-geth core.Genesis that
// op-chain-ops builds for it.
func fixture(f forkTimes) (*Genesis, *core.Genesis) {
	g := &Genesis{
		Config:        opChainConfig(f),
		Nonce:         7,
		Timestamp:     1234,
		ExtraData:     eip1559.EncodeJovianExtraData(250, 6, 0),
		GasLimit:      30_000_000,
		Difficulty:    big.NewInt(0),
		Mixhash:       common.Hash{0xaa},
		Coinbase:      predeploys.SequencerFeeVaultAddr,
		Alloc:         fixtureAlloc(),
		BaseFee:       big.NewInt(1_000_000_000),
		ExcessBlobGas: ptr.New[uint64](0),
		BlobGasUsed:   ptr.New[uint64](0),
	}
	gg := &core.Genesis{
		Config:        opGethChainConfig(f),
		Nonce:         g.Nonce,
		Timestamp:     g.Timestamp,
		ExtraData:     g.ExtraData,
		GasLimit:      g.GasLimit,
		Difficulty:    g.Difficulty,
		Mixhash:       g.Mixhash,
		Coinbase:      g.Coinbase,
		Alloc:         fixtureAlloc(),
		BaseFee:       g.BaseFee,
		ExcessBlobGas: g.ExcessBlobGas,
		BlobGasUsed:   g.BlobGasUsed,
	}
	return g, gg
}

// requireGenesisEqual compares big.Int fields by value: a decoded zero and big.NewInt(0)
// differ in their internal representation.
func requireGenesisEqual(t *testing.T, want, got *Genesis) {
	t.Helper()
	bigEq := cmp.Comparer(func(a, b *big.Int) bool {
		if a == nil || b == nil {
			return a == b
		}
		return a.Cmp(b) == 0
	})
	require.Empty(t, cmp.Diff(want, got, bigEq))
}

func parseJSON(t *testing.T, data []byte) map[string]any {
	t.Helper()
	var obj map[string]any
	require.NoError(t, json.Unmarshal(data, &obj))
	return obj
}

// TestMarshalMatchesOpGeth pins the encoding to the genesis.json that op-chain-ops writes
// through op-geth's core.Genesis. The only difference is osakaTime: GethChainConfig
// activates Osaka with Karst, which op-chain-ops' op-geth config omits.
func TestMarshalMatchesOpGeth(t *testing.T) {
	for _, tc := range []struct {
		name      string
		forks     forkTimes
		osakaTime any
	}{
		{name: "pre-Karst", forks: staggeredForks},
		{name: "Karst", forks: withKarst(staggeredForks, 9000), osakaTime: float64(9000)},
		{name: "Lagoon", forks: withLagoon(withKarst(staggeredForks, 9000), 10000), osakaTime: float64(9000)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			g, gg := fixture(tc.forks)
			g.Number, gg.Number = 3, 3
			g.GasUsed, gg.GasUsed = 21_000, 21_000
			g.ParentHash, gg.ParentHash = common.Hash{0xbb}, common.Hash{0xbb}
			g.SlotNumber, gg.SlotNumber = ptr.New[uint64](9), ptr.New[uint64](9)
			opGethJSON, err := json.Marshal(gg)
			require.NoError(t, err)
			want := parseJSON(t, opGethJSON)
			if tc.osakaTime != nil {
				want["config"].(map[string]any)["osakaTime"] = tc.osakaTime
			}

			got, err := json.Marshal(g)
			require.NoError(t, err)
			require.Equal(t, want, parseJSON(t, got))
		})
	}
}

// TestUnmarshalOpGeth decodes the genesis.json that op-geth's core.Genesis encodes.
func TestUnmarshalOpGeth(t *testing.T) {
	want, gg := fixture(withLagoon(withKarst(staggeredForks, 9000), 10000))
	data, err := json.Marshal(gg)
	require.NoError(t, err)

	var got Genesis
	require.NoError(t, json.Unmarshal(data, &got))
	requireGenesisEqual(t, want, &got)
}

func TestRoundTrip(t *testing.T) {
	g, _ := fixture(withLagoon(withKarst(staggeredForks, 9000), 10000))
	g.Number = 3
	g.GasUsed = 21_000
	g.ParentHash = common.Hash{0xbb}
	g.SlotNumber = ptr.New[uint64](9)

	data, err := json.Marshal(g)
	require.NoError(t, err)
	var got Genesis
	require.NoError(t, json.Unmarshal(data, &got))
	requireGenesisEqual(t, g, &got)

	again, err := json.Marshal(&got)
	require.NoError(t, err)
	require.JSONEq(t, string(data), string(again))
}

func TestEncoding(t *testing.T) {
	g, _ := fixture(staggeredForks)
	data, err := json.Marshal(g)
	require.NoError(t, err)
	obj := parseJSON(t, data)

	require.Equal(t, "0x7", obj["nonce"])
	require.Equal(t, "0x4d2", obj["timestamp"])
	require.Equal(t, "0x1c9c380", obj["gasLimit"])
	require.Equal(t, "0x0", obj["difficulty"])
	require.Equal(t, "0x3b9aca00", obj["baseFeePerGas"])
	require.Equal(t, common.Hash{0xaa}.Hex(), obj["mixHash"])
	require.Contains(t, obj["alloc"], "00000000000000000000000000000000000beef0", "alloc keys carry no 0x prefix")

	config := obj["config"].(map[string]any)
	require.Equal(t, float64(901), config["chainId"])
	require.Equal(t, float64(2000), config["canyonTime"])
	require.Equal(t, float64(2000), config["shanghaiTime"])
	require.Equal(t, float64(7000), config["pragueTime"])
	require.Equal(t, map[string]any{
		"eip1559Elasticity":        float64(6),
		"eip1559Denominator":       float64(50),
		"eip1559DenominatorCanyon": float64(250),
	}, config["optimism"])
}

func TestNilConfig(t *testing.T) {
	g, _ := fixture(staggeredForks)
	g.Config = nil
	data, err := json.Marshal(g)
	require.NoError(t, err)
	require.Nil(t, parseJSON(t, data)["config"])

	var got Genesis
	require.NoError(t, json.Unmarshal(data, &got))
	require.Nil(t, got.Config)
	require.Nil(t, got.GethGenesis().Config)
}

func TestUnmarshalMinimal(t *testing.T) {
	var got Genesis
	require.NoError(t, json.Unmarshal([]byte(`{"gasLimit":"0x1c9c380","difficulty":"0x0","alloc":{}}`), &got))
	requireGenesisEqual(t, &Genesis{GasLimit: 30_000_000, Difficulty: big.NewInt(0), Alloc: types.GenesisAlloc{}}, &got)
}

// TestUnmarshalConfigKeys checks how decoding treats "config" keys that Genesis does not
// store: unknown keys are ignored, and a derived key must match its derivation.
func TestUnmarshalConfigKeys(t *testing.T) {
	g, _ := fixture(staggeredForks)
	data, err := json.Marshal(g)
	require.NoError(t, err)

	decodeWithConfig := func(edit func(config map[string]any)) error {
		obj := parseJSON(t, data)
		edit(obj["config"].(map[string]any))
		edited, err := json.Marshal(obj)
		require.NoError(t, err)
		var got Genesis
		return json.Unmarshal(edited, &got)
	}

	require.NoError(t, decodeWithConfig(func(c map[string]any) { c["terminalTotalDifficultyPassed"] = true }),
		"unknown keys are ignored")
	require.NoError(t, decodeWithConfig(func(c map[string]any) { delete(c, "shanghaiTime") }),
		"derived keys may be absent")
	require.ErrorContains(t, decodeWithConfig(func(c map[string]any) { c["shanghaiTime"] = 1 }), `"shanghaiTime"`)
	require.ErrorContains(t, decodeWithConfig(func(c map[string]any) { c["londonBlock"] = 5 }), `"londonBlock"`)
}

func TestUnmarshalRejects(t *testing.T) {
	g, _ := fixture(staggeredForks)
	data, err := json.Marshal(g)
	require.NoError(t, err)

	for _, tc := range []struct {
		name   string
		edit   func(map[string]any)
		errMsg string
	}{
		{name: "missing gasLimit", edit: func(o map[string]any) { delete(o, "gasLimit") }, errMsg: "'gasLimit'"},
		{name: "missing difficulty", edit: func(o map[string]any) { delete(o, "difficulty") }, errMsg: "'difficulty'"},
		{name: "missing alloc", edit: func(o map[string]any) { delete(o, "alloc") }, errMsg: "'alloc'"},
		{name: "stateHash", edit: func(o map[string]any) { o["stateHash"] = common.Hash{0x01}.Hex() }, errMsg: "stateHash"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			obj := parseJSON(t, data)
			tc.edit(obj)
			edited, err := json.Marshal(obj)
			require.NoError(t, err)
			var got Genesis
			require.ErrorContains(t, json.Unmarshal(edited, &got), tc.errMsg)
		})
	}
}

// TestGethGenesisToBlock checks that the genesis block computed through GethGenesis is
// the one op-geth computes from the genesis that op-chain-ops builds, including the
// Isthmus withdrawals root (the L2ToL1MessagePasser storage root).
func TestGethGenesisToBlock(t *testing.T) {
	for _, tc := range []struct {
		name  string
		forks forkTimes
	}{
		{name: "pre-Isthmus", forks: staggeredForks},
		{name: "Isthmus", forks: isthmusAtGenesis},
		{name: "Karst", forks: withKarst(isthmusAtGenesis, 0)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			g, gg := fixture(tc.forks)
			want := gg.ToBlock()
			got := g.GethGenesis().ToBlock()
			require.Equal(t, want.Header(), got.Header())
			require.Equal(t, want.Hash(), got.Hash())
			if g.Config.IsIsthmus(g.Timestamp) {
				require.NotEqual(t, types.EmptyWithdrawalsHash, *got.Header().WithdrawalsHash,
					"Isthmus withdrawals root is the message passer's storage root")
			}
		})
	}
}
