package keeper

import (
	"context"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/certikfoundation/shentu/x/cert/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the cert MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) ProposeCertifier(goCtx context.Context, msg *types.MsgProposeCertifier) (*types.MsgProposeCertifierResponse, error) {
	return &types.MsgProposeCertifierResponse{}, nil
}

func (k msgServer) CertifyValidator(goCtx context.Context, msg *types.MsgCertifyValidator) (*types.MsgCertifyValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valPubKey, ok := msg.Pubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "Expecting cryptotypes.PubKey, got %T", valPubKey)
	}

	certifierAddr, err := sdk.AccAddressFromBech32(msg.Certifier)
	if err != nil {
		return nil, err
	}

	if err := k.Keeper.CertifyValidator(ctx, valPubKey, certifierAddr); err != nil {
		return nil, err
	}

	return &types.MsgCertifyValidatorResponse{}, nil
}

func (k msgServer) DecertifyValidator(goCtx context.Context, msg *types.MsgDecertifyValidator) (*types.MsgDecertifyValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valPubKey, ok := msg.Pubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "Expecting cryptotypes.PubKey, got %T", valPubKey)
	}

	decertifierAddr, err := sdk.AccAddressFromBech32(msg.Decertifier)
	if err != nil {
		return nil, err
	}

	if err := k.Keeper.DecertifyValidator(ctx, valPubKey, decertifierAddr); err != nil {
		return nil, err
	}

	return &types.MsgDecertifyValidatorResponse{}, nil
}

func (k msgServer) CertifyGeneral(goCtx context.Context, msg *types.MsgCertifyGeneral) (*types.MsgCertifyGeneralResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	certifierAddr, err := sdk.AccAddressFromBech32(msg.Certifier)
	if err != nil {
		return nil, err
	}

	certificate, err := types.NewGeneralCertificate(msg.CertificateType, msg.RequestContentType, msg.RequestContent, msg.Description, certifierAddr)
	if err != nil {
		return nil, err
	}

	certificateID, err := k.IssueCertificate(ctx, certificate)
	if err != nil {
		return nil, err
	}
	certEvent := sdk.NewEvent(
		types.EventTypeCertify,
		sdk.NewAttribute("certificate_id", certificateID.String()),
		sdk.NewAttribute("certificate_type", msg.CertificateType),
		sdk.NewAttribute("request_content_type", msg.RequestContentType),
		sdk.NewAttribute("request_content", msg.RequestContent),
		sdk.NewAttribute("description", msg.Description),
		sdk.NewAttribute("certifier", msg.Certifier),
	)
	ctx.EventManager().EmitEvent(certEvent)

	return &types.MsgCertifyGeneralResponse{}, nil
}

func (k msgServer) RevokeCertificate(goCtx context.Context, msg *types.MsgRevokeCertificate) (*types.MsgRevokeCertificateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	certificate, err := k.Keeper.GetCertificateByID(ctx, msg.Id)
	if err != nil {
		return nil, err
	}

	revokerAddr, err := sdk.AccAddressFromBech32(msg.Revoker)
	if err != nil {
		panic(err)
	}

	if err := k.Keeper.RevokeCertificate(ctx, certificate, revokerAddr); err != nil {
		return nil, err
	}
	revokeEvent := sdk.NewEvent(
		types.EventTypeRevokeCertificate,
		sdk.NewAttribute("revoker", msg.Revoker),
		sdk.NewAttribute("revoked_certificate", certificate.String()),
		sdk.NewAttribute("revoke_description", msg.Description),
	)
	ctx.EventManager().EmitEvent(revokeEvent)

	return &types.MsgRevokeCertificateResponse{}, nil
}

func (k msgServer) CertifyCompilation(goCtx context.Context, msg *types.MsgCertifyCompilation) (*types.MsgCertifyCompilationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	certifierAddr, err := sdk.AccAddressFromBech32(msg.Certifier)
	if err != nil {
		panic(err)
	}

	certificate := types.NewCompilationCertificate(
		types.CertificateTypeCompilation,
		msg.SourceCodeHash,
		msg.Compiler,
		msg.BytecodeHash,
		msg.Description,
		certifierAddr,
	)
	certificateID, err := k.Keeper.IssueCertificate(ctx, certificate)
	if err != nil {
		return nil, err
	}

	certEvent := sdk.NewEvent(
		types.EventTypeCertifyCompilation,
		sdk.NewAttribute("certificate_id", certificateID.String()),
		sdk.NewAttribute("source_code_hash", msg.SourceCodeHash),
		sdk.NewAttribute("compiler", msg.Compiler),
		sdk.NewAttribute("bytecode_hash", msg.BytecodeHash),
		sdk.NewAttribute("certifier", msg.Certifier),
	)
	ctx.EventManager().EmitEvent(certEvent)

	return &types.MsgCertifyCompilationResponse{}, nil
}

func (k msgServer) CertifyPlatform(goCtx context.Context, msg *types.MsgCertifyPlatform) (*types.MsgCertifyPlatformResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valPubKey, ok := msg.ValidatorPubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "Expecting cryptotypes.PubKey, got %T", valPubKey)
	}

	certifierAddr, err := sdk.AccAddressFromBech32(msg.Certifier)
	if err != nil {
		return nil, err
	}

	if err := k.Keeper.CertifyPlatform(ctx, certifierAddr, valPubKey, msg.Platform); err != nil {
		return nil, err
	}

	return &types.MsgCertifyPlatformResponse{}, nil
}
