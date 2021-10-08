package gql

import (
	"context"
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auctiontypes "github.com/tharsis/ethermint/x/auction/types"
	bondtypes "github.com/tharsis/ethermint/x/bond/types"
	nstypes "github.com/tharsis/ethermint/x/nameservice/types"
	"reflect"
	"strconv"
)

// OwnerAttributeName denotes the owner attribute name for a bond.
const OwnerAttributeName = "owner"

// BondIDAttributeName denotes the record bond ID.
const BondIDAttributeName = "bondId"

// ExpiryTimeAttributeName denotes the record expiry time.
const ExpiryTimeAttributeName = "expiryTime"

func getGQLCoin(coin sdk.Coin) *Coin {
	gqlCoin := Coin{
		Type:     coin.Denom,
		Quantity: strconv.FormatInt(coin.Amount.Int64(), 10),
	}

	return &gqlCoin
}

func getGQLCoins(coins sdk.Coins) []*Coin {
	gqlCoins := make([]*Coin, len(coins))
	for index, coin := range coins {
		gqlCoins[index] = getGQLCoin(coin)
	}

	return gqlCoins
}
func GetGQLNameAuthorityRecord(record *nstypes.NameAuthority) (*AuthorityRecord, error) {
	if record == nil {
		return nil, nil
	}

	return &AuthorityRecord{
		OwnerAddress:   record.OwnerAddress,
		OwnerPublicKey: record.OwnerPublicKey,
		Height:         strconv.FormatUint(record.Height, 10),
		Status:         record.Status,
		BondID:         record.GetBondId(),
		ExpiryTime:     record.GetExpiryTime().String(),
	}, nil
}

func getGQLRecord(ctx context.Context, resolver QueryResolver, record nstypes.Record) (*Record, error) {
	// Nil record.
	if record.Deleted {
		return nil, nil
	}

	recordType := record.ToRecordType()
	attributes, err := getAttributes(&recordType)
	if err != nil {
		return nil, err
	}

	references, err := getReferences(ctx, resolver, &recordType)
	if err != nil {
		return nil, err
	}

	return &Record{
		ID:         record.Id,
		BondID:     record.GetBondId(),
		CreateTime: record.GetCreateTime().String(),
		ExpiryTime: record.GetExpiryTime().String(),
		Owners:     record.GetOwners(),
		Attributes: attributes,
		References: references,
	}, nil
}

func getGQLNameRecord(record *nstypes.NameRecord) (*NameRecord, error) {
	if record == nil {
		return nil, nil
	}

	records := make([]*NameRecordEntry, len(record.History))
	for index, entry := range record.History {
		records[index] = getNameRecordEntry(entry)
	}

	return &NameRecord{
		Latest:  getNameRecordEntry(record.Latest),
		History: records,
	}, nil
}

func getNameRecordEntry(record *nstypes.NameRecordEntry) *NameRecordEntry {
	return &NameRecordEntry{
		ID:     record.Id,
		Height: strconv.FormatUint(record.Height, 10),
	}
}

func getGQLBond(bondObj *bondtypes.Bond) (*Bond, error) {
	// Nil record.
	if bondObj == nil {
		return nil, nil
	}

	return &Bond{
		ID:      bondObj.Id,
		Owner:   bondObj.Owner,
		Balance: getGQLCoins(bondObj.Balance),
	}, nil
}

func matchBondOnAttributes(bondObj *bondtypes.Bond, attributes []*KeyValueInput) bool {
	for _, attr := range attributes {
		switch attr.Key {
		case OwnerAttributeName:
			{
				if attr.Value.String == nil || bondObj.Owner != *attr.Value.String {
					return false
				}
			}
		default:
			{
				// Only attributes explicitly listed in the switch are queryable.
				return false
			}
		}
	}

	return true
}

func getAuctionBid(bid *auctiontypes.Bid) *AuctionBid {
	return &AuctionBid{
		BidderAddress: bid.BidderAddress,
		Status:        bid.Status,
		CommitHash:    bid.CommitHash,
		CommitTime:    bid.GetCommitTime(),
		RevealTime:    bid.GetRevealTime(),
		CommitFee:     getGQLCoin(bid.CommitFee),
		RevealFee:     getGQLCoin(bid.RevealFee),
		BidAmount:     getGQLCoin(bid.BidAmount),
	}
}

func GetGQLAuction(auction *auctiontypes.Auction, bids []*auctiontypes.Bid) (*Auction, error) {
	if auction == nil {
		return nil, nil
	}

	gqlAuction := Auction{
		ID:             auction.Id,
		Status:         auction.Status,
		OwnerAddress:   auction.OwnerAddress,
		CreateTime:     auction.GetCreateTime(),
		CommitsEndTime: auction.GetCommitsEndTime(),
		RevealsEndTime: auction.GetRevealsEndTime(),
		CommitFee:      getGQLCoin(auction.CommitFee),
		RevealFee:      getGQLCoin(auction.RevealFee),
		MinimumBid:     getGQLCoin(auction.MinimumBid),
		WinnerAddress:  auction.WinnerAddress,
		WinnerBid:      getGQLCoin(auction.WinningBid),
		WinnerPrice:    getGQLCoin(auction.WinningPrice),
	}

	auctionBids := make([]*AuctionBid, len(bids))
	for index, entry := range bids {
		auctionBids[index] = getAuctionBid(entry)
	}

	gqlAuction.Bids = auctionBids

	return &gqlAuction, nil
}

func getReferences(ctx context.Context, resolver QueryResolver, r *nstypes.RecordType) ([]*Record, error) {
	var ids []string

	for _, value := range r.Attributes {
		switch value.(type) {
		case interface{}:
			if obj, ok := value.(map[string]interface{}); ok {
				if _, ok := obj["/"]; ok && len(obj) == 1 {
					if _, ok := obj["/"].(string); ok {
						ids = append(ids, obj["/"].(string))
					}
				}
			}
		}
	}

	return resolver.GetRecordsByIds(ctx, ids)
}

func getAttributes(r *nstypes.RecordType) ([]*KeyValue, error) {
	return mapToKeyValuePairs(r.Attributes)
}

func mapToKeyValuePairs(attrs map[string]interface{}) ([]*KeyValue, error) {
	var kvPairs []*KeyValue

	trueVal := true
	falseVal := false

	for key, value := range attrs {
		kvPair := &KeyValue{
			Key:   key,
			Value: &Value{},
		}

		switch val := value.(type) {
		case nil:
			kvPair.Value.Null = &trueVal
		case int:
			kvPair.Value.Int = &val
		case float64:
			kvPair.Value.Float = &val
		case string:
			kvPair.Value.String = &val
		case bool:
			kvPair.Value.Boolean = &val
		case interface{}:
			if obj, ok := value.(map[string]interface{}); ok {
				if _, ok := obj["/"]; ok && len(obj) == 1 {
					if _, ok := obj["/"].(string); ok {
						kvPair.Value.Reference = &Reference{
							ID: obj["/"].(string),
						}
					}
				} else {
					bytes, err := json.Marshal(obj)
					if err != nil {
						return nil, err
					}

					jsonStr := string(bytes)
					kvPair.Value.JSON = &jsonStr
				}
			}
		}

		if kvPair.Value.Null == nil {
			kvPair.Value.Null = &falseVal
		}

		valueType := reflect.ValueOf(value)
		if valueType.Kind() == reflect.Slice {
			bytes, err := json.Marshal(value)
			if err != nil {
				return nil, err
			}

			jsonStr := string(bytes)
			kvPair.Value.JSON = &jsonStr
		}

		kvPairs = append(kvPairs, kvPair)
	}

	return kvPairs, nil
}
