package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	auctionkeeper "github.com/tharsis/ethermint/x/auction/keeper"
	bondkeeper "github.com/tharsis/ethermint/x/bond/keeper"
	"github.com/tharsis/ethermint/x/nameservice/helpers"
	"github.com/tharsis/ethermint/x/nameservice/types"
	"sort"
	"time"
)

var (

	// PrefixCIDToRecordIndex is the prefix for CID -> Record index.
	// Note: This is the primary index in the system.
	// Note: Golang doesn't support const arrays.
	PrefixCIDToRecordIndex = []byte{0x00}

	// PrefixNameAuthorityRecordIndex is the prefix for the name -> NameAuthority index.
	PrefixNameAuthorityRecordIndex = []byte{0x01}

	// PrefixWRNToNameRecordIndex is the prefix for the WRN -> NamingRecord index.
	PrefixWRNToNameRecordIndex = []byte{0x02}

	// PrefixBondIDToRecordsIndex is the prefix for the Bond ID -> [Record] index.
	PrefixBondIDToRecordsIndex = []byte{0x03}

	// PrefixBlockChangesetIndex is the prefix for the block changeset index.
	PrefixBlockChangesetIndex = []byte{0x04}

	// PrefixAuctionToAuthorityNameIndex is the prefix for the auction ID -> authority name index.
	PrefixAuctionToAuthorityNameIndex = []byte{0x05}

	// PrefixBondIDToAuthoritiesIndex is the prefix for the Bond ID -> [Authority] index.
	PrefixBondIDToAuthoritiesIndex = []byte{0x06}

	// PrefixExpiryTimeToRecordsIndex is the prefix for the Expiry Time -> [Record] index.
	PrefixExpiryTimeToRecordsIndex = []byte{0x10}

	// PrefixExpiryTimeToAuthoritiesIndex is the prefix for the Expiry Time -> [Authority] index.
	PrefixExpiryTimeToAuthoritiesIndex = []byte{0x11}

	// PrefixCIDToNamesIndex the the reverse index for naming, i.e. maps CID -> []Names.
	// TODO(ashwin): Move out of WNS once we have an indexing service.
	PrefixCIDToNamesIndex = []byte{0xe0}
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	accountKeeper auth.AccountKeeper
	bankKeeper    bank.Keeper
	recordKeeper  RecordKeeper
	bondKeeper    bondkeeper.Keeper
	auctionKeeper auctionkeeper.Keeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc codec.BinaryCodec // The wire codec for binary encoding/decoding.

	paramSubspace paramtypes.Subspace
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(cdc codec.BinaryCodec, accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper, recordKeeper RecordKeeper,
	bondKeeper bondkeeper.Keeper, auctionKeeper auctionkeeper.Keeper, storeKey sdk.StoreKey, ps paramtypes.Subspace) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}
	return Keeper{
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		recordKeeper:  recordKeeper,
		bondKeeper:    bondKeeper,
		auctionKeeper: auctionKeeper,
		storeKey:      storeKey,
		cdc:           cdc,
		paramSubspace: ps,
	}
}

// GetRecordIndexKey Generates Bond ID -> Bond index key.
func GetRecordIndexKey(id string) []byte {
	return append(PrefixCIDToRecordIndex, []byte(id)...)
}

// HasRecord - checks if a record by the given ID exists.
func (k Keeper) HasRecord(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(GetRecordIndexKey(id))
}

// GetRecord - gets a record from the store.
func (k Keeper) GetRecord(ctx sdk.Context, id string) (record types.Record) {
	store := ctx.KVStore(k.storeKey)
	result := store.Get(GetRecordIndexKey(id))
	k.cdc.MustUnmarshal(result, &record)
	return record
}

// ListRecords - get all records.
func (k Keeper) ListRecords(ctx sdk.Context) []types.Record {
	var records []types.Record

	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, PrefixCIDToRecordIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var obj types.Record
			k.cdc.MustUnmarshal(bz, &obj)
			//records = append(records, recordObjToRecord(store, k.cdc, obj))
			records = append(records, obj)
		}
	}

	return records
}

func (k Keeper) GetRecordExpiryQueue(ctx sdk.Context) []*types.ExpiryQueueRecord {
	var records []*types.ExpiryQueueRecord

	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, PrefixExpiryTimeToRecordsIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		record, err := helpers.BytesArrToStringArr(itr.Value())
		if err != nil {
			return records
		}
		records = append(records, &types.ExpiryQueueRecord{
			Id:    string(itr.Key()[len(PrefixExpiryTimeToRecordsIndex):]),
			Value: record,
		})
	}

	return records
}

// ProcessSetRecord creates a record.
func (k Keeper) ProcessSetRecord(ctx sdk.Context, msg types.MsgSetRecord) error {
	payload := msg.Payload.ToReadablePayload()
	record := types.RecordType{Attributes: payload.Record, BondId: msg.BondId}

	// Check signatures.
	resourceSignBytes, _ := record.GetSignBytes()
	cid, err := record.GetCID()
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid record JSON")
	}

	record.Id = cid

	if exists := k.HasRecord(ctx, record.Id); exists {
		return nil
	}

	record.Owners = []string{}
	for _, sig := range payload.Signatures {
		pubKey, err := legacy.PubKeyFromBytes(helpers.BytesFromBase64(sig.PubKey))
		if err != nil {
			fmt.Println("Error decoding pubKey from bytes: ", err)
			return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Invalid public key.")
		}

		sigOK := pubKey.VerifySignature(resourceSignBytes, helpers.BytesFromBase64(sig.Sig))
		if !sigOK {
			fmt.Println("Signature mismatch: ", sig.PubKey)
			return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Invalid signature.")
		}
		record.Owners = append(record.Owners, pubKey.Address().String())
	}

	// Sort owners list.
	sort.Strings(record.Owners)
	sdkErr := k.processRecord(ctx, &record, false)
	if sdkErr != nil {
		return sdkErr
	}
	return nil
}

func (k Keeper) processRecord(ctx sdk.Context, record *types.RecordType, isRenewal bool) error {
	params := k.GetParams(ctx)
	rent := params.RecordRent

	err := k.bondKeeper.TransferCoinsToModuleAccount(ctx, record.BondId, types.RecordRentModuleAccountName, sdk.NewCoins(rent))
	if err != nil {
		return err
	}

	record.CreateTime = ctx.BlockHeader().Time
	record.ExpiryTime = ctx.BlockHeader().Time.Add(params.RecordRentDuration)
	record.Deleted = false

	k.PutRecord(ctx, record.ToRecordObj())
	k.InsertRecordExpiryQueue(ctx, record.ToRecordObj())

	// Renewal doesn't change the name and bond indexes.
	if !isRenewal {
		k.AddBondToRecordIndexEntry(ctx, record.BondId, record.Id)
	}

	return nil
}

// PutRecord - saves a record to the store and updates ID -> Record index.
func (k Keeper) PutRecord(ctx sdk.Context, record types.Record) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetRecordIndexKey(record.Id), k.cdc.MustMarshal(&record))
	k.updateBlockChangeSetForRecord(ctx, record.Id)
}

// AddBondToRecordIndexEntry adds the Bond ID -> [Record] index entry.
func (k Keeper) AddBondToRecordIndexEntry(ctx sdk.Context, bondID string, id string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(getBondIDToRecordsIndexKey(bondID, id), []byte{})
}

// Generates Bond ID -> Records index key.
func getBondIDToRecordsIndexKey(bondID string, id string) []byte {
	return append(append(PrefixBondIDToRecordsIndex, []byte(bondID)...), []byte(id)...)
}

// getRecordExpiryQueueTimeKey gets the prefix for the record expiry queue.
func getRecordExpiryQueueTimeKey(timestamp time.Time) []byte {
	timeBytes := sdk.FormatTimeBytes(timestamp)
	return append(PrefixExpiryTimeToRecordsIndex, timeBytes...)
}

// SetRecordExpiryQueueTimeSlice sets a specific record expiry queue timeslice.
func (k Keeper) SetRecordExpiryQueueTimeSlice(ctx sdk.Context, timestamp time.Time, cids []string) {
	store := ctx.KVStore(k.storeKey)
	bz, _ := helpers.StrArrToBytesArr(cids)
	store.Set(getRecordExpiryQueueTimeKey(timestamp), bz)
}

// DeleteRecordExpiryQueueTimeSlice deletes a specific record expiry queue timeslice.
func (k Keeper) DeleteRecordExpiryQueueTimeSlice(ctx sdk.Context, timestamp time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(getRecordExpiryQueueTimeKey(timestamp))
}

// GetRecordExpiryQueueTimeSlice gets a specific record queue timeslice.
// A timeslice is a slice of CIDs corresponding to records that expire at a certain time.
func (k Keeper) GetRecordExpiryQueueTimeSlice(ctx sdk.Context, timestamp time.Time) []string {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(getRecordExpiryQueueTimeKey(timestamp))
	if bz == nil {
		return []string{}
	}
	cids, err := helpers.BytesArrToStringArr(bz)
	if err != nil {
		return []string{}
	}
	return cids
}

// InsertRecordExpiryQueue inserts a record CID to the appropriate timeslice in the record expiry queue.
func (k Keeper) InsertRecordExpiryQueue(ctx sdk.Context, val types.Record) {
	timeSlice := k.GetRecordExpiryQueueTimeSlice(ctx, val.ExpiryTime)
	timeSlice = append(timeSlice, val.Id)
	k.SetRecordExpiryQueueTimeSlice(ctx, val.ExpiryTime, timeSlice)
}

// GetModuleBalances gets the nameservice module account(s) balances.
func (k Keeper) GetModuleBalances(ctx sdk.Context) []*types.AccountBalance {
	var balances []*types.AccountBalance
	accountNames := []string{types.RecordRentModuleAccountName, types.AuthorityRentModuleAccountName}

	for _, accountName := range accountNames {
		moduleAddress := k.accountKeeper.GetModuleAddress(accountName)
		moduleAccount := k.accountKeeper.GetAccount(ctx, moduleAddress)
		if moduleAccount != nil {
			accountBalance := k.bankKeeper.GetAllBalances(ctx, moduleAddress)
			balances = append(balances, &types.AccountBalance{
				AccountName: accountName,
				Balance:     accountBalance,
			})
		}
	}

	return balances
}
