// Package keeper specifies the keeper for the cert module.
package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/cert/types"
)

// Keeper manages certifier & security council related logics.
type Keeper struct {
	storeKey       sdk.StoreKey
	cdc            codec.BinaryMarshaler
	slashingKeeper types.SlashingKeeper
	stakingKeeper  types.StakingKeeper
}

// NewKeeper creates a new instance of the certifier keeper.
func NewKeeper(cdc codec.BinaryMarshaler, storeKey sdk.StoreKey, slashingKeeper types.SlashingKeeper, stakingKeeper types.StakingKeeper) Keeper {
	return Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		slashingKeeper: slashingKeeper,
		stakingKeeper:  stakingKeeper,
	}
}

// CertifyPlatform certifies a validator host platform by a certifier.
func (k Keeper) CertifyPlatform(ctx sdk.Context, certifier sdk.AccAddress, validator cryptotypes.PubKey, description string) error {
	if !k.IsCertifier(ctx, certifier) {
		return types.ErrRejectedValidator
	}

	pkAny, err := codectypes.NewAnyWithValue(validator)
	if err != nil {
		return err
	}

	bz := k.cdc.MustMarshalBinaryBare(&types.Platform{ValidatorPubkey: pkAny, Description: description})
	ctx.KVStore(k.storeKey).Set(types.PlatformStoreKey(validator), bz)
	return nil
}

// GetPlatform returns the host platform of the validator.
func (k Keeper) GetPlatform(ctx sdk.Context, validator cryptotypes.PubKey) (string, bool) {
	if platform := ctx.KVStore(k.storeKey).Get(types.PlatformStoreKey(validator)); platform != nil {
		return string(platform), true
	}
	return "", false
}

// GetAllPlatforms gets all platform certificates for genesis export
func (k Keeper) GetAllPlatforms(ctx sdk.Context) (platforms []types.Platform) {
	platforms = make([]types.Platform, 0)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PlatformsStoreKey())
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var platform types.Platform
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &platform)
		platforms = append(platforms, platform)
	}
	return platforms
}
