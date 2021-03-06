package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/x/cert/keeper"
	"github.com/certikfoundation/shentu/x/cert/types"
)

const (
	OpWeightMsgCertifyValidator = "op_weight_msg_certify_validator"
	OpWeightMsgCertifyPlatform  = "op_weight_msg_certify_platform"
	OpWeightMsgCertifyAuditing  = "op_weight_msg_certify_auditing"
	OpWeightMsgCertifyProof     = "op_weight_msg_certify_proof"
)

// Default simulation operation weights for messages.
const (
	DefaultWeightMsgCertify int = 20
)

// WeightedOperations creates an operation (with weight) for each type of message generators.
func WeightedOperations(appParams simtypes.AppParams, cdc codec.JSONMarshaler, ak types.AccountKeeper,
	bk types.BankKeeper, k keeper.Keeper) simulation.WeightedOperations {
	var weightMsgCertifyValidator int
	appParams.GetOrGenerate(cdc, OpWeightMsgCertifyValidator, &weightMsgCertifyValidator, nil,
		func(_ *rand.Rand) {
			weightMsgCertifyValidator = simappparams.DefaultWeightMsgSend
		})

	var weightMsgCertifyPlatform int
	appParams.GetOrGenerate(cdc, OpWeightMsgCertifyPlatform, &weightMsgCertifyPlatform, nil,
		func(_ *rand.Rand) {
			weightMsgCertifyPlatform = simappparams.DefaultWeightMsgSend
		})

	var weightMsgCertifyAuditing int
	appParams.GetOrGenerate(cdc, OpWeightMsgCertifyAuditing, &weightMsgCertifyAuditing, nil,
		func(_ *rand.Rand) {
			weightMsgCertifyAuditing = simappparams.DefaultWeightMsgSend
		})

	var weightMsgCertifyProof int
	appParams.GetOrGenerate(cdc, OpWeightMsgCertifyProof, &weightMsgCertifyProof, nil,
		func(_ *rand.Rand) {
			weightMsgCertifyProof = simappparams.DefaultWeightMsgSend
		})

	var weightMsgCertifyIdentity int
	appParams.GetOrGenerate(cdc, OpWeightMsgCertifyProof, &weightMsgCertifyIdentity, nil,
		func(_ *rand.Rand) {
			weightMsgCertifyIdentity = DefaultWeightMsgCertify
		})

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(weightMsgCertifyValidator, SimulateMsgCertifyValidator(ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgCertifyPlatform, SimulateMsgCertifyPlatform(ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgCertifyAuditing, SimulateMsgCertifyAuditing(ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgCertifyProof, SimulateMsgCertifyProof(ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgCertifyIdentity, SimulateMsgCertifyIdentity(ak, bk, k)),
	}
}

// SimulateMsgCertifyValidator generates a MsgCertifyValidator object which fields contain
// a randomly chosen existing certifier and randomized validator's PubKey.
func SimulateMsgCertifyValidator(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		certifiers := k.GetAllCertifiers(ctx)
		certifier := certifiers[r.Intn(len(certifiers))]
		certifierAddr, err := sdk.AccAddressFromBech32(certifier.Address)
		if err != nil {
			panic(err)
		}
		var certifierAcc simtypes.Account
		for _, acc := range accs {
			if acc.Address.Equals(certifierAddr) {
				certifierAcc = acc
				break
			}
		}
		validator := simtypes.RandomAccounts(r, 1)[0]

		msg, err := types.NewMsgCertifyValidator(certifierAddr, validator.PubKey)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCertifyValidator, err.Error()), nil, err
		}

		account := ak.GetAccount(ctx, certifierAddr)
		fees, err := simtypes.RandomFees(r, ctx, bk.SpendableCoins(ctx, account.GetAddress()))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCertifyValidator, err.Error()), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			certifierAcc.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCertifyValidator, err.Error()), nil, err
		}
		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgCertifyPlatform generates a MsgCertifyPlatform object which fields contain
// a randomly chosen existing certifier, a randomized validator's PubKey and a random string description.
func SimulateMsgCertifyPlatform(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		certifiers := k.GetAllCertifiers(ctx)
		certifier := certifiers[r.Intn(len(certifiers))]
		certifierAddr, err := sdk.AccAddressFromBech32(certifier.Address)
		if err != nil {
			panic(err)
		}
		var certifierAcc simtypes.Account
		for _, acc := range accs {
			if acc.Address.Equals(certifierAddr) {
				certifierAcc = acc
				break
			}
		}
		validator := simtypes.RandomAccounts(r, 1)[0]
		platform := simtypes.RandStringOfLength(r, 10)

		msg, err := types.NewMsgCertifyPlatform(certifierAddr, validator.PubKey, platform)
		if err != nil {
			panic(err)
		}

		account := ak.GetAccount(ctx, certifierAddr)
		fees, err := simtypes.RandomFees(r, ctx, bk.SpendableCoins(ctx, account.GetAddress()))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCertifyPlatform, err.Error()), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			certifierAcc.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}
		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgCertifyAuditing generates a MsgCertifyAuditing object which fields contain
// a randomly chosen existing certifer, a random contract and a random string description.
func SimulateMsgCertifyAuditing(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account,
		chainID string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		certifiers := k.GetAllCertifiers(ctx)
		certifier := certifiers[r.Intn(len(certifiers))]
		certifierAddr, err := sdk.AccAddressFromBech32(certifier.Address)
		if err != nil {
			panic(err)
		}
		var certifierAcc simtypes.Account
		for _, acc := range accs {
			if acc.Address.Equals(certifierAddr) {
				certifierAcc = acc
				break
			}
		}
		contract := simtypes.RandomAccounts(r, 1)[0]
		description := simtypes.RandStringOfLength(r, 10)

		msg := types.NewMsgCertifyGeneral("auditing", "address", contract.Address.String(), description, certifierAddr)

		account := ak.GetAccount(ctx, certifierAddr)
		fees, err := simtypes.RandomFees(r, ctx, bk.SpendableCoins(ctx, account.GetAddress()))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			certifierAcc.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}
		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgCertifyProof generates a MsgCertifyProof object which fields contain
// a randomly chosen existing certifer, a random contract and a random string description.
func SimulateMsgCertifyProof(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		certifiers := k.GetAllCertifiers(ctx)
		certifier := certifiers[r.Intn(len(certifiers))]
		certifierAddr, err := sdk.AccAddressFromBech32(certifier.Address)
		if err != nil {
			panic(err)
		}
		var certifierAcc simtypes.Account
		for _, acc := range accs {
			if acc.Address.Equals(certifierAddr) {
				certifierAcc = acc
				break
			}
		}
		contract := simtypes.RandomAccounts(r, 1)[0]
		description := simtypes.RandStringOfLength(r, 10)

		msg := types.NewMsgCertifyGeneral("proof", "address", contract.Address.String(), description, certifierAddr)

		account := ak.GetAccount(ctx, certifierAddr)
		fees, err := simtypes.RandomFees(r, ctx, bk.SpendableCoins(ctx, account.GetAddress()))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			certifierAcc.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}
		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgCertifyIdentity generates a MsgCertifyGeneral object to certify a random account address.
func SimulateMsgCertifyIdentity(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		certifiers := k.GetAllCertifiers(ctx)
		certifier := certifiers[r.Intn(len(certifiers))]
		certifierAddr, err := sdk.AccAddressFromBech32(certifier.Address)
		if err != nil {
			panic(err)
		}

		var certifierAcc simtypes.Account
		for _, acc := range accs {
			if acc.Address.Equals(certifierAddr) {
				certifierAcc = acc
				break
			}
		}

		delAddr, found := keeper.RandomDelegator(r, k, ctx)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCertifyGeneral, err.Error()), nil, nil
		}
		identityAcc := ak.GetAccount(ctx, delAddr)

		msg := types.NewMsgCertifyGeneral("identity", "address", identityAcc.GetAddress().String(), "", certifierAddr)

		account := ak.GetAccount(ctx, certifierAddr)
		fees, err := simtypes.RandomFees(r, ctx, bk.SpendableCoins(ctx, account.GetAddress()))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			certifierAcc.PrivKey,
		)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}
		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
