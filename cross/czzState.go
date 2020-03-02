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


type EntangleState struct {
	EnInfos 		map[czzutil.Address]*LightHouseInfo
	EnEntitys 		map[uint64]UserEntangleInfos
	CurExchangeID 	uint64
}

/////////////////////////////////////////////////////////////////
// keep staking enough amount asset
func (es *EntangleState) RegisterLightHouse(addr czzutil.Address,amount *big.Int,
	fee uint64,assetType uint32) error {
	if amount.Cmp(MinStakingAmountForLightHouse) < 0 {
		return ErrLessThanMin
	}
	if _,ok := es.EnInfos[addr]; ok {
		return ErrRepeatRegister
	}
	info := &LightHouseInfo{
		ExchangeID:		es.CurExchangeID+1,
		Address:		addr,
		StakingAmount:	new(big.Int).Set(amount),
		AssetFlag:		assetType,
		Fee:			fee,
		EnAssets:		make([]*EnAssetItem,0,0),
		EntangleAmount:	big.NewInt(0),
		WhiteList:		make([]*WhiteUnit,0,0),
	}
	es.EnInfos[addr] = info
	return nil
}
func (es *EntangleState) AppendWhiteList(addr czzutil.Address,wlist []*WhiteUnit) error {
	if val,ok := es.EnInfos[addr]; ok {
		cnt := len(val.WhiteList)
		if cnt + len(wlist) >= MaxWhiteListCount {
			return errors.New("more than max white list")
		}
		for _,v := range wlist {
			if ValidAssetType(v.AssetType) && ValidPK(v.Pk) {
				val.WhiteList = append(val.WhiteList,v)
			}
		}
		return nil
	} else {
		return ErrNoRegister
	}	
}
// UnregisterLightHouse need to check all the proves and handle all the user's burn coins
func (es *EntangleState) UnregisterLightHouse(addr czzutil.Address) error {
	if val,ok := es.EnInfos[addr]; ok {
		last := new(big.Int).Sub(val.StakingAmount,val.EntangleAmount)
		redeemAmount(addr,last)
	} else {
		return ErrNoRegister
	}
	return nil
}
// AddEntangleItem add item in the state, keep lighthouse have enough amount to entangle,
func (es *EntangleState) AddEntangleItem(addr czzutil.Address,aType uint32,lightID uint64,
	height,amount *big.Int) (*big.Int,error) {
	lh := es.getLightHouse(lightID)
	if lh == nil {
		return nil,ErrNoRegister
	}
	if !isValidAsset(aType,lh.AssetFlag) {
		return nil,ErrNoUserAsset
	}
	sendAmount := big.NewInt(0)
	var err error
	lhEntitys,ok := es.EnEntitys[lightID]
	if !ok {
		lhEntitys = UserEntangleInfos(make(map[czzutil.Address]EntangleEntitys))
	}
	if lhEntitys != nil {
		userEntitys,ok1 := lhEntitys[addr]
		if !ok1 {
			userEntitys = EntangleEntitys(make([]*EntangleEntity,0,0))
		}
		found := false
		for _,v := range userEntitys {
			if aType == v.AssetType {
				found = true
				v.EntangleAmount = new(big.Int).Add(v.EntangleAmount,amount)
				break
			}
		}
		if !found {
			entity := &EntangleEntity{
				ExchangeID:	lightID,
				Address:	addr,
				AssetType:	aType,
				Height:		new(big.Int).Set(height),
				EntangleAmount: new(big.Int).Set(amount),
				BurnAmount: newBurnInfos(),
			}
			userEntitys = append(userEntitys,entity)
		}
		
		// calc the send amount
		reserve := es.getEntangledAmount(lightID,aType)
		sendAmount,err = calcEntangleAmount(reserve,amount,aType)
		if err != nil {
			return nil,err
		}
		lh.addEnAsset(aType,amount)
		lh.recordEntangleAmount(sendAmount)
		lhEntitys[addr] = userEntitys
		es.EnEntitys[lightID] = lhEntitys
	}
	return sendAmount,nil
}
// BurnAsset user burn the czz asset to exchange the outside asset,the caller keep the burn was true.
// verify the txid,keep equal amount czz
func (es *EntangleState) BurnAsset(addr czzutil.Address,aType uint32,lightID uint64,
	amount *big.Int) (*big.Int,error) {
	light := es.getLightHouse(lightID)
	if light == nil {
		return nil,ErrNoRegister
	}
	lhEntitys,ok := es.EnEntitys[lightID]
	if !ok {
		return nil,ErrNoRegister
	}
	userEntitys,ok1 := lhEntitys[addr]
	if !ok1 {
		return nil,ErrNoUserReg
	}
	var userEntity *EntangleEntity
	for _,v := range userEntitys {
		if aType == v.AssetType {
			userEntity = v
			break
		}
	}
	if userEntity == nil {
		return nil,ErrNoUserAsset
	}
	validAmount := userEntity.BurnAmount.GetValidAmount()
	if amount.Cmp(validAmount) > 0 {
		return nil,ErrNotEnouthBurn
	} 
	userEntity.BurnAmount.BurnAmount = new(big.Int).Add(userEntity.BurnAmount.BurnAmount,amount)
	res := new(big.Int).Div(new(big.Int).Mul(amount,big.NewInt(int64(light.Fee))),big.NewInt(int64(light.Fee)))
	
	return res,nil
}
func (es *EntangleState) ConfiscateAsset() error {
	return nil
}
//////////////////////////////////////////////////////////////////////
func redeemAmount(addr czzutil.Address,amount *big.Int) error {
	if amount.Sign() > 0 {
		
	}
	return nil
}
func calcEntangleAmount(reserve,reqAmount *big.Int,atype uint32) (*big.Int,error) {
	return nil,nil
}
func (es *EntangleState) getEntangledAmount(lightID uint64,atype uint32) *big.Int {
	aa := big.NewInt(0)
	if lhEntitys,ok := es.EnEntitys[lightID];ok {
		for _,userEntitys := range lhEntitys{
			for _,vv := range userEntitys {
				if atype == vv.AssetType {
					aa = aa.Add(aa,vv.EntangleAmount)
					break
				}
			}
		}
	}
	return aa
}
func (es *EntangleState) getLightHouse(id uint64) *LightHouseInfo {
	for _,val := range es.EnInfos {
		if val.ExchangeID == id {
			return val
		}
	}
	return nil
}
func (es *EntangleState) getAllEntangleAmount(atype uint32) *big.Int {
	all := big.NewInt(0)
	for _,val := range es.EnInfos {
		for _,v := range val.EnAssets {
			if v.AssetType == atype {
				all = all.Add(all,v.Amount)
				break
			}
		}
	}
	return all
}
// 最低质押额度＝ 100 万 CZZ ＋（累计跨链买入 CZZ －累计跨链卖出 CZZ）x 汇率比
func (es *EntangleState) LimitStakingAmount(eid uint64,atype uint32) *big.Int {
	lh := es.getLightHouse(eid)
	if lh != nil {
		l := new(big.Int).Sub(lh.StakingAmount,lh.EntangleAmount)
		if l.Sign() > 0 {
			l = new(big.Int).Sub(l,MinStakingAmountForLightHouse)
			if l.Sign() > 0 {
				return l
			}
		}
	}
	return nil
}

//////////////////////////////////////////////////////////////////////
