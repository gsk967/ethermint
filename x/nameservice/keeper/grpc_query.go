package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tharsis/ethermint/x/nameservice/types"
)

type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

func (q Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := q.Keeper.GetParams(ctx)
	return &types.QueryParamsResponse{Params: &params}, nil
}

func (q Querier) ListRecords(c context.Context, _ *types.QueryListRecordsRequest) (*types.QueryListRecordsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	records := q.Keeper.ListRecords(ctx)
	return &types.QueryListRecordsResponse{Records: records}, nil
}

func (q Querier) GetRecord(c context.Context, req *types.QueryRecordByIdRequest) (*types.QueryRecordByIdResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	id := req.GetId()
	if !q.Keeper.HasRecord(ctx, id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Record not found.")
	}
	record := q.Keeper.GetRecord(ctx, id)
	return &types.QueryRecordByIdResponse{Record: record}, nil
}

func (q Querier) GetRecordByBondId(c context.Context, req *types.QueryRecordByBondIdRequest) (*types.QueryRecordByBondIdResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	records := q.recordKeeper.QueryRecordsByBond(ctx, req.GetId())
	return &types.QueryRecordByBondIdResponse{Records: records}, nil
}

func (q Querier) GetNameServiceModuleBalance(c context.Context, _ *types.GetNameServiceModuleBalanceRequest) (*types.GetNameServiceModuleBalanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	balances := q.Keeper.GetModuleBalances(ctx)
	return &types.GetNameServiceModuleBalanceResponse{
		Balances: balances,
	}, nil
}

func (q Querier) ListNameRecords(c context.Context, _ *types.QueryListNameRecordsRequest) (*types.QueryListNameRecordsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	nameRecords := q.Keeper.ListNameRecords(ctx)
	return &types.QueryListNameRecordsResponse{Names: nameRecords}, nil
}

func (q Querier) Whois(c context.Context, request *types.QueryWhoisRequest) (*types.QueryWhoisResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	nameAuthority := q.Keeper.GetNameAuthority(ctx, request.GetName())
	return &types.QueryWhoisResponse{NameAuthority: nameAuthority}, nil
}

func (q Querier) LookupWrn(c context.Context, req *types.QueryLookupWrn) (*types.QueryLookupWrnResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	wrn := req.GetWrn()
	if !q.Keeper.HasNameRecord(ctx, wrn) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "WRN not found.")
	}
	nameRecord := q.Keeper.GetNameRecord(ctx, wrn)
	if nameRecord == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "name record not found.")
	}
	return &types.QueryLookupWrnResponse{Name: nameRecord}, nil
}

func (q Querier) ResolveWrn(c context.Context, req *types.QueryResolveWrn) (*types.QueryResolveWrnResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	wrn := req.GetWrn()
	record := q.Keeper.ResolveWRN(ctx, wrn)
	if record == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "record not found.")
	}
	return &types.QueryResolveWrnResponse{Record: record}, nil
}

func (q Querier) GetRecordExpiryQueue(c context.Context, _ *types.QueryGetRecordExpiryQueue) (*types.QueryGetRecordExpiryQueueResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	records := q.Keeper.GetRecordExpiryQueue(ctx)
	return &types.QueryGetRecordExpiryQueueResponse{Records: records}, nil
}

func (q Querier) GetAuthorityExpiryQueue(c context.Context, _ *types.QueryGetAuthorityExpiryQueue) (*types.QueryGetAuthorityExpiryQueueResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	authorities := q.Keeper.GetAuthorityExpiryQueue(ctx)
	return &types.QueryGetAuthorityExpiryQueueResponse{Authorities: authorities}, nil
}
