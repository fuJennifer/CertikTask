package cert

import (
	"github.com/gogo/protobuf/proto"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/certikfoundation/shentu/x/cert/keeper"
	"github.com/certikfoundation/shentu/x/cert/types"
)

func InitDefaultGenesis(ctx sdk.Context, k keeper.Keeper) {
	InitGenesis(ctx, k, *types.DefaultGenesisState())
}

// InitGenesis initialize default parameters and the keeper's address to pubkey map.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	certifiers := data.Certifiers
	validators := data.Validators
	platforms := data.Platforms
	certificates := data.Certificates
	libraries := data.Libraries

	for _, certifier := range certifiers {
		k.SetCertifier(ctx, certifier)
	}
	if len(certifiers) > 0 {
		certifierAddr, err := sdk.AccAddressFromBech32(certifiers[0].Address)
		if err != nil {
			panic(err)
		}
		for _, platform := range platforms {
			pk, ok := platform.ValidatorPubkey.GetCachedValue().(cryptotypes.PubKey)
			if !ok {
				panic(sdkerrors.Wrapf(sdkerrors.ErrUnpackAny, "cannot unpack Any into cryto.PubKey %T", platform.ValidatorPubkey))
			}

			_ = k.CertifyPlatform(ctx, certifierAddr, pk, platform.Description)
		}
	}
	for _, validator := range validators {
		pk, ok := validator.Pubkey.GetCachedValue().(cryptotypes.PubKey)
		if !ok {
			panic(sdkerrors.Wrapf(sdkerrors.ErrUnpackAny, "cannot unpack Any into cryto.PubKey %T", validator.Pubkey))
		}
		certifierAddr, err := sdk.AccAddressFromBech32(validator.Certifier)
		if err != nil {
			panic(err)
		}

		k.SetValidator(ctx, pk, certifierAddr)
	}
	for _, certificateAny := range certificates {
		certificate, ok := certificateAny.GetCachedValue().(types.Certificate)
		if !ok {
			panic(sdkerrors.Wrapf(sdkerrors.ErrUnpackAny, "cannot unpack Any into Certificate %T", certificateAny))
		}
		k.SetCertificate(ctx, certificate)
	}
	for _, library := range libraries {
		libAddr, err := sdk.AccAddressFromBech32(library.Address)
		if err != nil {
			panic(err)
		}
		publisherAddr, err := sdk.AccAddressFromBech32(library.Publisher)
		if err != nil {
			panic(err)
		}
		k.SetLibrary(ctx, libAddr, publisherAddr)
	}
}

// ExportGenesis writes the current store values to a genesis file, which can be imported again with InitGenesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	certifiers := k.GetAllCertifiers(ctx)
	validators := k.GetAllValidators(ctx)
	platforms := k.GetAllPlatforms(ctx)
	certificates := k.GetAllCertificates(ctx)
	libraries := k.GetAllLibraries(ctx)

	certificateAnys := make([]*codectypes.Any, len(certificates))
	for i, certificate := range certificates {
		msg, ok := certificate.(proto.Message)
		if !ok {
			panic(sdkerrors.Wrapf(sdkerrors.ErrPackAny, "cannot proto marshal %T", certificate))
		}
		any, err := codectypes.NewAnyWithValue(msg)
		if err != nil {
			panic(err)
		}
		certificateAnys[i] = any
	}

	return &types.GenesisState{
		Certifiers:   certifiers,
		Validators:   validators,
		Platforms:    platforms,
		Certificates: certificateAnys,
		Libraries:    libraries,
	}
}
