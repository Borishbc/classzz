package cross

import (
	// "bytes"
	// "encoding/binary"
	"errors"
	// "fmt"
	// "strings"
	"math/big"

	// "github.com/classzz/classzz/chaincfg"
	// "github.com/classzz/classzz/chaincfg/chainhash"
	// "github.com/classzz/classzz/czzec"
	// "github.com/classzz/classzz/txscript"
	// "github.com/classzz/classzz/wire"
	"github.com/classzz/czzutil"
)

var (
	ErrInvalidParam      = errors.New("Invalid Param")
	ErrLessThanMin		 = errors.New("less than min staking amount for lighthouse")
	ErrRepeatRegister    = errors.New("repeat register on this address")
	ErrNoRegister 	     = errors.New("not found the lighthouse")
	ErrNoUserReg 	     = errors.New("not entangle user in the lighthouse")
	ErrNoUserAsset 	     = errors.New("user no entangle asset in the lighthouse")
	ErrNotEnouthBurn 	 = errors.New("not enough burn amount in lighthouse")
)

var (
	MinStakingAmountForLightHouse = new(big.Int).Mul(big.NewInt(1000000),big.NewInt(1e9))
	MaxWhiteListCount  = 5
	MAXBASEFEE 		   = 10000
)
const (
	LhAssetBTC uint32 = 1 << iota
	LhAssetBCH
	LhAssetBSV
	LhAssetLTC
	LhAssetUSDT
	LhAssetDOGE
)
type BurnItem struct {
	Amount 	*big.Int			// czz asset amount
	Height 	uint64	
}
type BurnInfos struct {
	Items 	[]*BurnItem
	BurnAmount *big.Int 		// burned amount for outside asset
}
func newBurnInfos() *BurnInfos {
	return nil
}
func (b *BurnInfos) GetAllAmount() *big.Int {
	amount := big.NewInt(0)
	for _,v := range b.Items {
		amount = amount.Add(amount,v.Amount)
	}
	return amount
}
func (b *BurnInfos) GetValidAmount() *big.Int {
	return nil
}
// Update the valid amount for diffence height for entangle info
func (b *BurnInfos) Update() {

}
type WhiteUnit struct {
	AssetType 		uint32
	Pk				[]byte
}
type BaseAmountUint struct {
	AssetType 		uint32
	Amount 			*big.Int
}
type EnAssetItem BaseAmountUint
type FreeQuotaItem BaseAmountUint

type LightHouseInfo struct {
	ExchangeID		uint64
	Address 		czzutil.Address
	StakingAmount 	*big.Int 			// in 
	EntangleAmount  *big.Int			// out,express by czz
	EnAssets		[]*EnAssetItem		// out,the extrinsic asset
	Frees 			[]*FreeQuotaItem	// extrinsic asset
	AssetFlag 		uint32
	Fee 			uint64
	KeepTime		uint64 		// the time as the block count for finally redeem time
	WhiteList 		[]*WhiteUnit
}

func (lh *LightHouseInfo) addEnAsset(atype uint32, amount *big.Int) {
	found := false
	for _,val := range lh.EnAssets {
		if val.AssetType == atype {
			found = true
			val.Amount = new(big.Int).Add(val.Amount,amount)			
		}
	}
	if !found {
		lh.EnAssets = append(lh.EnAssets,&EnAssetItem{
			AssetType:		atype,
			Amount:			amount,
		})
	}
}
func (lh *LightHouseInfo) recordEntangleAmount(amount *big.Int) {
	lh.EntangleAmount = new(big.Int).Add(lh.EntangleAmount,amount)
} 
func (lh *LightHouseInfo) addFreeQuota(amount *big.Int,atype uint32) {
	for _,v := range lh.Frees {
		if atype == v.AssetType {
			v.Amount = new(big.Int).Add(v.Amount,amount)
		}
	}
}
func (lh *LightHouseInfo) useFreeQuota(amount *big.Int,atype uint32) {
	for _,v := range lh.Frees {
		if atype == v.AssetType {
			if v.Amount.Cmp(amount) >= 0 {
				v.Amount = new(big.Int).Sub(v.Amount,amount)
			} else {
				// panic
				v.Amount = big.NewInt(0)
			}
		}
	}
}
func (lh *LightHouseInfo) canRedeem(amount *big.Int,atype uint32) bool {
	for _,v := range lh.Frees {
		if atype == v.AssetType {
			if v.Amount.Cmp(amount) >= 0 {
				return true
			} else {
				return false
			}
		}
	}
	return false
}
/////////////////////////////////////////////////////////////////
// Address > EntangleEntity
type EntangleEntity struct {
	ExchangeID		uint64
	Address 		czzutil.Address
	AssetType		uint32
	Height			*big.Int				// newest height for entangle
	OldHeight		*big.Int				// oldest height for entangle
	EntangleAmount 	*big.Int				// out asset	
	MaxRedeem       *big.Int				// out asset
	BurnAmount 		*BurnInfos
}
type EntangleEntitys []*EntangleEntity
type UserEntangleInfos map[czzutil.Address]EntangleEntitys


/////////////////////////////////////////////////////////////////
func (e *EntangleEntity) updateFreeQuota(limit *big.Int) {
	
}
func (e *EntangleEntity) GetValidRedeemAmount() *big.Int {
	return e.MaxRedeem
}

func (ee *EntangleEntitys) updateFreeQuotaForAllType(limit big.Int) *big.Int {
	return nil
}

/////////////////////////////////////////////////////////////////

func ValidAssetType(utype uint32) bool {
	if utype & LhAssetBTC != 0 || utype & LhAssetBCH != 0 || utype & LhAssetBSV != 0 ||
	utype & LhAssetLTC != 0 || utype & LhAssetUSDT != 0 || utype & LhAssetDOGE != 0 {
		return true
	}
	return false
}
func ValidPK(pk []byte) bool {
	return true
}
func isValidAsset(atype,assetAll uint32) bool {
	return atype & assetAll != 0
}