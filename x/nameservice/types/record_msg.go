package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgSetRecord{}
	_ sdk.Msg = &MsgRenewRecord{}
	_ sdk.Msg = &MsgAssociateBond{}
	_ sdk.Msg = &MsgDissociateBond{}
	_ sdk.Msg = &MsgDissociateRecords{}
	_ sdk.Msg = &MsgReAssociateRecords{}
)

// NewMsgSetRecord is the constructor function for MsgSetRecord.
func NewMsgSetRecord(payload Payload, bondID string, signer sdk.AccAddress) MsgSetRecord {
	return MsgSetRecord{
		Payload: payload,
		BondId:  bondID,
		Signer:  signer.String(),
	}
}

func (msg MsgSetRecord) ValidateBasic() error {
	if len(msg.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer)
	}
	owners := msg.Payload.Record.Owners
	for _, owner := range owners {
		if owner == "" {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Record owner not set.")
		}
	}

	if len(msg.BondId) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Bond ID is required.")
	}
	return nil
}

func (msg MsgSetRecord) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{accAddr}
}

// GetSignBytes gets the sign bytes for the msg MsgCreateBond
func (msg MsgSetRecord) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// NewMsgRenewRecord is the constructor function for MsgRenewRecord.
func NewMsgRenewRecord(recordId string, signer sdk.AccAddress) MsgRenewRecord {
	return MsgRenewRecord{
		RecordId: recordId,
		Signer:   signer.String(),
	}
}

// Route Implements Msg.
func (msg MsgRenewRecord) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgRenewRecord) Type() string { return "renew-record" }

// ValidateBasic Implements Msg.
func (msg MsgRenewRecord) ValidateBasic() error {
	if len(msg.RecordId) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "record id is required.")
	}

	if len(msg.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer.")
	}

	return nil
}

// GetSignBytes gets the sign bytes for Msg
func (msg MsgRenewRecord) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgRenewRecord) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{accAddr}
}

// NewMsgAssociateBond is the constructor function for MsgAssociateBond.
func NewMsgAssociateBond(recordId, bondId string, signer sdk.AccAddress) MsgAssociateBond {
	return MsgAssociateBond{
		BondId:   bondId,
		RecordId: recordId,
		Signer:   signer.String(),
	}
}

// Route Implements Msg.
func (msg MsgAssociateBond) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgAssociateBond) Type() string { return "associate-bond" }

// ValidateBasic Implements Msg.
func (msg MsgAssociateBond) ValidateBasic() error {
	if len(msg.RecordId) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "record id is required.")
	}
	if len(msg.BondId) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "bond id is required.")
	}
	if len(msg.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer.")
	}

	return nil
}

// GetSignBytes gets the sign bytes for Msg
func (msg MsgAssociateBond) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgAssociateBond) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{accAddr}
}

// NewMsgDissociateBond is the constructor function for MsgDissociateBond.
func NewMsgDissociateBond(recordId string, signer sdk.AccAddress) MsgDissociateBond {
	return MsgDissociateBond{
		RecordId: recordId,
		Signer:   signer.String(),
	}
}

// Route Implements Msg.
func (msg MsgDissociateBond) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgDissociateBond) Type() string { return "dissociate-bond" }

// ValidateBasic Implements Msg.
func (msg MsgDissociateBond) ValidateBasic() error {
	if len(msg.RecordId) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "record id is required.")
	}
	if len(msg.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer.")
	}

	return nil
}

// GetSignBytes gets the sign bytes for Msg
func (msg MsgDissociateBond) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgDissociateBond) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{accAddr}
}

// NewMsgDissociateRecords is the constructor function for MsgDissociateRecords.
func NewMsgDissociateRecords(bondId string, signer sdk.AccAddress) MsgDissociateRecords {
	return MsgDissociateRecords{
		BondId: bondId,
		Signer: signer.String(),
	}
}

// Route Implements Msg.
func (msg MsgDissociateRecords) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgDissociateRecords) Type() string { return "dissociate-records" }

// ValidateBasic Implements Msg.
func (msg MsgDissociateRecords) ValidateBasic() error {
	if len(msg.BondId) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "bond id is required.")
	}
	if len(msg.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer.")
	}

	return nil
}

// GetSignBytes gets the sign bytes for Msg
func (msg MsgDissociateRecords) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgDissociateRecords) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{accAddr}
}

// NewMsgReAssociateRecords is the constructor function for MsgReAssociateRecords.
func NewMsgReAssociateRecords(oldBondId, newBondId string, signer sdk.AccAddress) MsgReAssociateRecords {
	return MsgReAssociateRecords{
		OldBondId: oldBondId,
		NewBondId: newBondId,
		Signer:    signer.String(),
	}
}

// Route Implements Msg.
func (msg MsgReAssociateRecords) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgReAssociateRecords) Type() string { return "reassociate-records" }

// ValidateBasic Implements Msg.
func (msg MsgReAssociateRecords) ValidateBasic() error {
	if len(msg.OldBondId) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "old-bond-id is required.")
	}
	if len(msg.NewBondId) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "new-bond-id is required.")
	}
	if len(msg.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer.")
	}

	return nil
}

// GetSignBytes gets the sign bytes for Msg
func (msg MsgReAssociateRecords) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgReAssociateRecords) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{accAddr}
}
